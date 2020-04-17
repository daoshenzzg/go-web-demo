package http

import (
	"bytes"
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go-web-demo/library/net/netutil/breaker"
	xtime "go-web-demo/library/time"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	_minRead     = 16 * 1024 // 16kb
	_contentType = "Content-Type"
	_urlencoded  = "application/x-www-form-urlencoded"
	_userAgent   = "User-Agent"
)

var (
	_noKickUserAgent = "zhiguang@mgtv.com"
)

func init() {
	n, err := os.Hostname()
	if err == nil {
		_noKickUserAgent = _noKickUserAgent + " " + runtime.Version() + " " + n
	}
}

// ClientConfig is http client conf.
type ClientConfig struct {
	MaxTotal    int
	MaxPerHost  int
	KeepAlive   xtime.Duration
	DialTimeout xtime.Duration
	Timeout     xtime.Duration
	Breaker     *breaker.Config
}

// Client is http client.
type Client struct {
	conf      *ClientConfig
	client    *http.Client
	dialer    *net.Dialer
	transport http.Transport
	mutex     sync.RWMutex
	breaker   *breaker.Group
}

// NewClient new a http client pool
func NewClient(c *ClientConfig) *Client {
	if c.DialTimeout <= 0 || c.Timeout <= 0 {
		panic("must config http timeout")
	}

	client := new(Client)
	client.breaker = breaker.NewGroup(c.Breaker)
	client.conf = c
	client.dialer = &net.Dialer{
		Timeout:   time.Duration(c.DialTimeout),
		KeepAlive: time.Duration(c.KeepAlive),
	}
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         client.dialer.DialContext,
		MaxIdleConns:        c.MaxTotal,
		MaxIdleConnsPerHost: c.MaxPerHost,
		IdleConnTimeout:     time.Duration(c.KeepAlive),
	}
	client.client = &http.Client{
		Transport: transport,
		Timeout:   time.Duration(c.Timeout),
	}
	return client
}

// SetTransport set client transport
func (client *Client) SetTransport(t http.Transport) {
	client.transport = t
}

// SetConfig set client config.
func (client *Client) SetConfig(c *ClientConfig) {
	client.mutex.Lock()
	if c.MaxTotal > 0 {
		client.conf.MaxTotal = c.MaxTotal
	}
	if c.MaxPerHost > 0 {
		client.conf.MaxPerHost = c.MaxPerHost
	}
	if c.Timeout > 0 {
		client.client.Timeout = time.Duration(c.Timeout)
		client.conf.Timeout = c.Timeout
	}
	if c.DialTimeout > 0 {
		client.dialer.Timeout = time.Duration(c.DialTimeout)
		client.conf.DialTimeout = c.DialTimeout
	}
	if c.KeepAlive > 0 {
		client.dialer.KeepAlive = time.Duration(c.KeepAlive)
		client.conf.KeepAlive = c.KeepAlive
	}
	if c.Breaker != nil {
		client.conf.Breaker = c.Breaker
		client.breaker.Reload(c.Breaker)
	}
	client.mutex.Unlock()
}

// NewRequest new http request with method, uri, values and headers.
func (client *Client) NewRequest(method, uri string, params url.Values) (req *http.Request, err error) {
	ru := uri
	if params != nil {
		ru = uri + "?" + params.Encode()
	}
	req, err = http.NewRequest(http.MethodGet, ru, nil)
	if err != nil {
		err = errors.Wrapf(err, "method:%s,uri:%s", method, ru)
		return
	}
	if method == http.MethodPost {
		req.Header.Set(_contentType, _urlencoded)
	}
	req.Header.Set(_userAgent, _noKickUserAgent)
	return
}

// Get issues a GET to the specified URL.
func (client *Client) Get(c context.Context, uri string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(http.MethodGet, uri, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// Post issues a Post to the specified URL.
func (client *Client) Post(c context.Context, uri string, params url.Values, res interface{}) (err error) {
	req, err := client.NewRequest(http.MethodPost, uri, params)
	if err != nil {
		return
	}
	return client.Do(c, req, res)
}

// JSON sends an HTTP request and returns an HTTP json response.
func (client *Client) JSON(c context.Context, req *http.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if err = json.Unmarshal(bs, res); err != nil {
			err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// Raw sends an HTTP request and returns bytes response
func (client *Client) Raw(c context.Context, req *http.Request, v ...string) (bs []byte, err error) {
	var (
		resp *http.Response
		uri  = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.Host, req.URL.Path)
	)

	// NOTE fix prom & config uri key.
	if len(v) == 1 {
		uri = v[0]
	}

	err = client.breaker.Do(uri, func() error {
		if resp, err = client.client.Do(req); err != nil {
			return errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
		defer resp.Body.Close()
		if resp.StatusCode >= http.StatusBadRequest {
			return errors.Errorf("incorrect http status:%d host:%s, url:%s", resp.StatusCode, req.URL.Host, realURL(req))
		}
		if bs, err = readAll(resp.Body, _minRead); err != nil {
			return errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
		return nil
	})
	return
}

// Do sends an HTTP request and returns an HTTP json response.
func (client *Client) Do(c context.Context, req *http.Request, res interface{}, v ...string) (err error) {
	var bs []byte
	if bs, err = client.Raw(c, req, v...); err != nil {
		return
	}
	if res != nil {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if err = json.Unmarshal(bs, res); err != nil {
			err = errors.Wrapf(err, "host:%s, url:%s", req.URL.Host, realURL(req))
		}
	}
	return
}

// realUrl return url with http://host/params.
func realURL(req *http.Request) string {
	if req.Method == http.MethodGet {
		return req.URL.String()
	} else if req.Method == http.MethodPost {
		ru := req.URL.Path
		if req.Body != nil {
			rd, ok := req.Body.(io.Reader)
			if ok {
				buf := bytes.NewBuffer([]byte{})
				buf.ReadFrom(rd)
				ru = ru + "?" + buf.String()
			}
		}
		return ru
	}
	return req.URL.Path
}

// readAll reads from r until an error or EOF and returns the data it read
// from the internal buffer allocated with a specified capacity.
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r)
	return buf.Bytes(), err
}
