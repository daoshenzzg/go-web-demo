package rpc

import (
	"context"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"gopkg.in/go-playground/assert.v1"
	"net"
	"sync"
	"testing"
)

func TestPool_Put(t *testing.T) {
	// setup server
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})

	go s.Serve(l)
	defer s.Stop()

	pool := NewPool(nil)
	if err := pool.DialAll(context.Background(), l.Addr().String()); err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer pool.Close()
	assert.Equal(t, pool.client.conf.PoolSize, len(pool.conns))
}

func TestPool_GetConn(t *testing.T) {
	// setup server
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})

	go s.Serve(l)
	defer s.Stop()

	pool := NewPool(nil)
	if err := pool.DialAll(context.Background(), l.Addr().String()); err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer pool.Close()

	wg := new(sync.WaitGroup)
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				pool.GetConn()
			}
		}()
	}
	wg.Wait()
}

func TestGetConnAndSayHello(t *testing.T) {
	// setup server
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})

	go s.Serve(l)
	defer s.Stop()

	pool := NewPool(nil)
	if err := pool.DialAll(context.Background(), l.Addr().String()); err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer pool.Close()

	c := pb.NewGreeterClient(pool.GetConn())
	rsp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "John"})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	if rsp.Message != "Hello John" {
		t.Fatalf("Got unexpected response %v", rsp.Message)
	}
	t.Log(rsp.Message)
}
