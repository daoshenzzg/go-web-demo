package rpc

import (
	"context"
	"google.golang.org/grpc"
	"math"
	"sync"
	"sync/atomic"
)

var (
	_abortSize uint64 = math.MaxUint32 / 2
)

type Pool struct {
	mu      sync.RWMutex
	counter uint64
	client  *Client
	conns   []*grpc.ClientConn
}

// NewPool NewPool
func NewPool(conf *ClientConfig, opt ...grpc.DialOption) *Pool {
	client := NewClient(conf, opt...)
	return &Pool{
		client: client,
		conns:  make([]*grpc.ClientConn, 0),
	}
}

// Put put a grpc ClientConn into Pool
func (p *Pool) Put(conn *grpc.ClientConn) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.conns)+1 > p.client.conf.PoolSize {
		return false
	}
	p.conns = append(p.conns, conn)
	return true
}

// GetConn get conns from pool
func (p *Pool) GetConn() *grpc.ClientConn {
	v := atomic.AddUint64(&p.counter, 1)
	if v > _abortSize {
		atomic.StoreUint64(&p.counter, 0)
	}
	idx := v % uint64(len(p.conns))
	return p.conns[idx]
}

// DialAll dial grpc conns
func (p *Pool) DialAll(ctx context.Context, target string, opt ...grpc.DialOption) (err error) {
	for i := 0; i < p.client.conf.PoolSize; i++ {
		conn, err := p.client.Dial(ctx, target, opt...)
		if err != nil {
			return err
		}
		p.Put(conn)
	}
	return
}

// Close close all pool conns
func (p *Pool) Close() {
	for _, conn := range p.conns {
		conn.Close()
	}
}
