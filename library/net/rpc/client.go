package rpc

import (
	"context"
	"github.com/pkg/errors"
	"go-web-demo/library/ecode"
	"go-web-demo/library/log"
	"go-web-demo/library/net/netutil/breaker"
	"go-web-demo/library/net/rpc/status"
	xtime "go-web-demo/library/time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	gstatus "google.golang.org/grpc/status"
	"math"
	"sync"
	"time"
)

var (
	_once        sync.Once
	_defaultConf = &ClientConfig{
		DialTimeout: xtime.Duration(10 * time.Second),
		PoolSize:    4,
		Timeout:     xtime.Duration(250 * time.Millisecond),
	}
	_defaultClient *Client
	_abortIndex    int8 = math.MaxInt8 / 2
)

type ClientConfig struct {
	DialTimeout xtime.Duration
	Timeout     xtime.Duration
	PoolSize    int
	NonBlock    bool
	Breaker     *breaker.Config
}

// Client is the framework's client side instance, it contains the ctx, opt and interceptors.
// Create an instance of Client, by using NewClient().
type Client struct {
	conf    *ClientConfig
	breaker *breaker.Group
	mutex   sync.RWMutex

	opt      []grpc.DialOption
	handlers []grpc.UnaryClientInterceptor
}

// handle returns a new unary client interceptor for OpenTracing\Logging\LinkTimeout.
func (c *Client) handle() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		var (
			cancel context.CancelFunc
			p      peer.Peer
		)
		var ec ecode.Codes = ecode.OK

		_, ctx, cancel = c.conf.Timeout.Shrink(ctx)
		defer cancel()

		err = c.breaker.Do(method, func() error {
			opts = append(opts, grpc.Peer(&p))
			err = invoker(ctx, method, req, reply, cc, opts...)
			if err != nil {
				gst, _ := gstatus.FromError(err)
				ec = status.ToECode(gst)
				err = errors.WithMessage(ec, gst.Message())
				log.Error("hystrix got error(%v)", err)
			}
			return err
		})
		return
	}
}

// SetConfig hot reloads client config
func (c *Client) SetConfig(conf *ClientConfig) (err error) {
	if conf == nil {
		conf = _defaultConf
	}
	if conf.DialTimeout <= 0 {
		conf.DialTimeout = _defaultConf.DialTimeout
	}
	if conf.Timeout <= 0 {
		conf.Timeout = _defaultConf.Timeout
	}
	if conf.PoolSize <= 0 {
		conf.PoolSize = _defaultConf.PoolSize
	}

	c.mutex.Lock()
	c.conf = conf
	if c.breaker == nil {
		c.breaker = breaker.NewGroup(conf.Breaker)
	} else {
		c.breaker.Reload(conf.Breaker)
	}
	c.mutex.Unlock()
	return nil
}

// Use attachs a global inteceptor to the Client.
// For example, this is the right place for a circuit breaker or error management inteceptor.
func (c *Client) Use(handlers ...grpc.UnaryClientInterceptor) *Client {
	finalSize := len(c.handlers) + len(handlers)
	if finalSize >= int(_abortIndex) {
		panic("rrpc: client use too many handlers")
	}
	mergedHandlers := make([]grpc.UnaryClientInterceptor, finalSize)
	copy(mergedHandlers, c.handlers)
	copy(mergedHandlers[len(c.handlers):], handlers)
	c.handlers = mergedHandlers
	return c
}

// UseOpt attachs a global rpc DialOption to the Client.
func (c *Client) UseOpt(opt ...grpc.DialOption) *Client {
	c.opt = append(c.opt, opt...)
	return c
}

// NewConn will create a rpc conns by default config.
func NewConn(target string, opt ...grpc.DialOption) (*grpc.ClientConn, error) {
	return DefaultClient().Dial(context.Background(), target, opt...)
}

// NewClient returns a new blank Client instance with a default client interceptor.
// opt can be used to add rpc dial options.
func NewClient(conf *ClientConfig, opt ...grpc.DialOption) *Client {
	c := new(Client)
	if err := c.SetConfig(conf); err != nil {
		panic(err)
	}
	c.UseOpt(opt...)
	c.Use(c.recovery(), c.handle())
	return c
}

// DefaultClient returns a new default Client instance with a default client interceptor and default dialoption.
// opt can be used to add rpc dial options.
func DefaultClient() *Client {
	_once.Do(func() {
		_defaultClient = NewClient(nil)
	})
	return _defaultClient
}

// Dial creates a client connection to the given target.
func (c *Client) Dial(ctx context.Context, target string, opt ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	if !c.conf.NonBlock {
		c.opt = append(c.opt, grpc.WithBlock())
	}
	c.opt = append(c.opt, grpc.WithInsecure())
	c.opt = append(c.opt, grpc.WithUnaryInterceptor(c.chainUnaryClient()))
	c.opt = append(c.opt, opt...)
	c.mutex.RLock()
	conf := c.conf
	c.mutex.RUnlock()
	if conf.DialTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(conf.DialTimeout))
		defer cancel()
	}
	if conn, err = grpc.DialContext(ctx, target, c.opt...); err != nil {
		log.Error("client: dial %s error %v!", target, err)
	}
	err = errors.WithStack(err)
	return
}

// chainUnaryClient creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func (c *Client) chainUnaryClient() grpc.UnaryClientInterceptor {
	n := len(c.handlers)
	if n == 0 {
		return func(ctx context.Context, method string, req, reply interface{},
			cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var (
			i            int
			chainHandler grpc.UnaryInvoker
		)
		chainHandler = func(ictx context.Context, imethod string, ireq, ireply interface{}, ic *grpc.ClientConn, iopts ...grpc.CallOption) error {
			if i == n-1 {
				return invoker(ictx, imethod, ireq, ireply, ic, iopts...)
			}
			i++
			return c.handlers[i](ictx, imethod, ireq, ireply, ic, chainHandler, iopts...)
		}

		return c.handlers[0](ctx, method, req, reply, cc, chainHandler, opts...)
	}
}
