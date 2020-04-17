package http

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-web-demo/library/log"
	"go-web-demo/library/net/netutil/breaker"
	"go-web-demo/library/render"
	xtime "go-web-demo/library/time"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	readTimeout := xtime.Duration(time.Second)
	writeTimeout := xtime.Duration(time.Second)
	endPoint := fmt.Sprintf(":%d", 8081)
	maxHeaderBytes := 1 << 20
	httpSrv := &http.Server{
		Addr:           endPoint,
		Handler:        engine,
		ReadTimeout:    time.Duration(readTimeout),
		WriteTimeout:   time.Duration(writeTimeout),
		MaxHeaderBytes: maxHeaderBytes,
	}
	engine.GET("/mytest", func(ctx *gin.Context) {
		time.Sleep(time.Millisecond * 500)
		r := render.New(ctx)
		r.JSON("", nil)
	})
	engine.GET("/mytest1", func(ctx *gin.Context) {
		time.Sleep(time.Millisecond * 500)
		r := render.New(ctx)
		r.JSON("", nil)
	})

	go func() {
		// service connections
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("srv.ListenAndServe() error(%v)", err)
			panic(err)
		}
	}()

	client := NewClient(
		&ClientConfig{
			MaxTotal:    20,
			MaxPerHost:  20,
			KeepAlive:   xtime.Duration(time.Second),
			DialTimeout: xtime.Duration(time.Second),
			Timeout:     xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Namespace:              "testGroup",
				Timeout:                1 * xtime.Duration(time.Second),
				MaxConcurrentRequests:  100,
				RequestVolumeThreshold: 10,
				SleepWindow:            5 * xtime.Duration(time.Second),
				ErrorPercentThreshold:  50,
			},
		})

	var res struct {
		Code int `json:"code"`
	}

	// test Get
	params := url.Values{}
	params.Set("type", "5")
	params.Set("clipId", "336434")
	params.Set("version", "5.5")

	if err := client.Get(context.Background(), "http://10.1.172.179:8101/odin/p1/coll/info", params, &res); err != nil {
		t.Errorf("HTTPClient: expected no error but got %v, res %v", err, res)
	}
	if res.Code != 200 {
		t.Errorf("HTTPClient: expected code=0 but got %d res %v", res.Code, res)
	}
	// test Post
	err := client.Post(context.Background(), "http://10.1.172.179:8101/odin/p1/coll/info", params, &res)
	if err != nil {
		t.Errorf("HTTPClient: expected no error but got %v", err)
	}
	// test server and timeout.
	client.SetConfig(&ClientConfig{KeepAlive: xtime.Duration(time.Second * 20), Timeout: xtime.Duration(time.Millisecond * 400)})
	if err := client.Get(context.Background(), "http://localhost:8081/mytest", nil, &res); err == nil {
		fmt.Printf("code %v", res.Code)
		t.Errorf("HTTPClient: expected error timeout for request")
	}
	client.SetConfig(&ClientConfig{Timeout: xtime.Duration(time.Millisecond * 300)})
	if err := client.Get(context.Background(), "http://10.1.172.179:8101/odin/p1/coll/info", params, &res); err != nil {
		t.Errorf("HTTPClient: expected no error but got %v", err)
	}
	client.SetConfig(&ClientConfig{Timeout: xtime.Duration(time.Millisecond * 1)})
	if err := client.Get(context.Background(), "http://10.1.172.179:8101/odin/p1/coll/info", params, &res); err == nil {
		t.Errorf("HTTPClient: expected error timeout but got %v", err)
	}
	client.SetConfig(&ClientConfig{KeepAlive: xtime.Duration(time.Second * 70)})
}

func TestDo(t *testing.T) {
	var (
		clipId  = 336434
		version = "5.5"
		uri     = "http://10.1.172.179:8101/odin/p1/coll/info"
		req     *http.Request
		client  *Client
		err     error
	)
	client = NewClient(
		&ClientConfig{
			DialTimeout: xtime.Duration(time.Second),
			Timeout:     xtime.Duration(time.Second),
			KeepAlive:   xtime.Duration(time.Second),
			Breaker: &breaker.Config{
				Namespace:              "testGroup",
				Timeout:                1 * xtime.Duration(time.Second),
				MaxConcurrentRequests:  50,
				RequestVolumeThreshold: 10,
				SleepWindow:            5 * xtime.Duration(time.Millisecond),
				ErrorPercentThreshold:  50,
			},
		})
	params := url.Values{}
	params.Set("clipId", strconv.Itoa(clipId))
	params.Set("version", version)
	if req, err = client.NewRequest("GET", uri, params); err != nil {
		t.Errorf("client.NewRequest: get error(%v)", err)
	}
	var res struct {
		Code int `json:"code"`
	}
	if err = client.Do(context.TODO(), req, &res); err != nil {
		t.Errorf("Do: client.Do get error(%v) url: %s", err, realURL(req))
	}
}

func BenchmarkDo(b *testing.B) {
	cf := &ClientConfig{
		DialTimeout: xtime.Duration(time.Second),
		Timeout:     xtime.Duration(time.Second),
		KeepAlive:   xtime.Duration(time.Second),
		Breaker: &breaker.Config{
			Namespace:              "testGroup",
			Timeout:                1 * xtime.Duration(time.Second),
			MaxConcurrentRequests:  10,
			RequestVolumeThreshold: 1,
			SleepWindow:            1 * xtime.Duration(time.Millisecond),
			ErrorPercentThreshold:  50,
		},
	}
	var (
		clipId  = 336434
		version = "5.5"
		uri     = "http://10.1.172.179:8101/odin/p1/coll/info"
		req     *http.Request
		client  *Client
		err     error
	)
	client = NewClient(cf)
	b.ResetTimer()
	b.N = 10
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// client.SetConfig(cf)
			params := url.Values{}
			params.Set("clipId", strconv.Itoa(clipId))
			params.Set("version", version)
			req, err = client.NewRequest(http.MethodGet, uri, params)
			if err != nil {
				b.Errorf("newRequest: get error(%v)", err)
				continue
			}
			var res struct {
				Code int `json:"code"`
			}
			if err = client.Do(context.TODO(), req, &res); err != nil {
				b.Errorf("Do: client.Do get error(%v)", err)
			}
		}
	})

	uri = "http://10.1.172.179:8101/odin/p1/coll/infox" // NOTE: for breaker
	b.ResetTimer()
	b.N = 10
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// client.SetConfig(cf)
			params := url.Values{}
			params.Set("clipId", strconv.Itoa(clipId))
			params.Set("version", version)
			req, err := client.NewRequest(http.MethodGet, uri, params)
			if err != nil {
				b.Errorf("newRequest: get error(%v)", err)
				continue
			}
			var res struct {
				Code int `json:"code"`
			}
			if err = client.Do(context.TODO(), req, &res); err != nil {
				b.Logf("Do: client.Do get error(%v)", err)
			}
		}
	})
}
