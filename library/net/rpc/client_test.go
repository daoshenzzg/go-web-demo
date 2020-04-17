package rpc

import (
	"context"
	"errors"
	"fmt"
	"go-web-demo/library/net/netutil/breaker"
	xtime "go-web-demo/library/time"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
	"net"
	"testing"
	"time"
)

// server is used to implement helloworld.GreeterServer.
type greeterServer struct{}

// SayHello implements helloworld.GreeterServer
func (g *greeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "Error" {
		return nil, errors.New("error")
	}
	if in.Name == "InvalidArgument" {
		return nil, status.Error(codes.InvalidArgument, "InvalidArgument")
	}
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func TestClient(t *testing.T) {
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

	// setup client
	client := NewClient(&ClientConfig{
		DialTimeout: xtime.Duration(time.Second * 10),
		Timeout:     xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Namespace:              "testGRPC",
			Timeout:                1 * xtime.Duration(time.Second),
			MaxConcurrentRequests:  1000,
			RequestVolumeThreshold: 10,
			SleepWindow:            5 * xtime.Duration(time.Second),
			ErrorPercentThreshold:  50,
		},
	})
	// apply client interceptor middleware
	client.Use(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (ret error) {
		newctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		ret = invoker(newctx, method, req, reply, cc, opts...)
		return
	})
	conn, err := client.Dial(context.Background(), l.Addr().String())
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	rsp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "John"})
	if err != nil {
		t.Fatalf("could not greet: %v", err)
	}
	if rsp.Message != "Hello John" {
		t.Fatalf("Got unexpected response %v", rsp.Message)
	}
	fmt.Println(rsp.Message)
}

func BenchmarkGrpc(b *testing.B) {
	// setup server
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		b.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})

	go s.Serve(l)
	defer s.Stop()

	// setup client
	client := NewClient(&ClientConfig{
		DialTimeout: xtime.Duration(time.Second * 10),
		Timeout:     xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Namespace:              "testGRPC",
			Timeout:                1 * xtime.Duration(time.Second),
			MaxConcurrentRequests:  100,
			RequestVolumeThreshold: 10,
			SleepWindow:            5 * xtime.Duration(time.Second),
			ErrorPercentThreshold:  50,
		},
	})
	// apply client interceptor middleware
	client.Use(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (ret error) {
		newctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		ret = invoker(newctx, method, req, reply, cc, opts...)
		return
	})
	conn, err := client.Dial(context.Background(), l.Addr().String())
	if err != nil {
		b.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	in := &pb.HelloRequest{Name: "John"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := c.SayHello(context.Background(), in)
			if err != nil {
				b.Logf("could not greet: %v", err)
			}
		}
	})
}

func BenchmarkGrpcBreaker(b *testing.B) {
	// setup server
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		b.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeterServer{})

	go s.Serve(l)
	defer s.Stop()

	// setup client
	client := NewClient(&ClientConfig{
		DialTimeout: xtime.Duration(time.Second * 10),
		Timeout:     xtime.Duration(time.Second * 10),
		Breaker: &breaker.Config{
			Namespace:              "testGRPC",
			Timeout:                1 * xtime.Duration(time.Second),
			MaxConcurrentRequests:  100,
			RequestVolumeThreshold: 10,
			SleepWindow:            5 * xtime.Duration(time.Second),
			ErrorPercentThreshold:  10,
		},
	})
	// apply client interceptor middleware
	client.Use(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (ret error) {
		newctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		ret = invoker(newctx, method, req, reply, cc, opts...)
		return
	})
	conn, err := client.Dial(context.Background(), l.Addr().String())
	if err != nil {
		b.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	in := &pb.HelloRequest{Name: "InvalidArgument"} // NOTE: for breaker
	b.ResetTimer()
	b.N = 10
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := c.SayHello(context.Background(), in); err != nil {
				b.Logf("c.SayHello get error(%v)", err)
			}
		}
	})
}
