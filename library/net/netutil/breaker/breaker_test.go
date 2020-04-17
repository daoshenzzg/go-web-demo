package breaker

import (
	"go-web-demo/library/ecode"
	xtime "go-web-demo/library/time"
	"github.com/pkg/errors"
	"sync"
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	g1 := NewGroup(nil)
	g2 := NewGroup(_conf)
	if g1.conf != g2.conf {
		t.FailNow()
	}

	g := NewGroup(_conf)
	c := &Config{
		Namespace:              "test",
		Timeout:                xtime.Duration(1 * time.Second),
		MaxConcurrentRequests:  100,
		RequestVolumeThreshold: 10,
		SleepWindow:            xtime.Duration(5 * time.Second),
		ErrorPercentThreshold:  50,
	}
	g.Reload(c)
	if g.conf.Namespace == _conf.Namespace {
		t.FailNow()
	}
}

func TestDo(t *testing.T) {
	c := &Config{
		Namespace:              "test",
		Timeout:                xtime.Duration(1 * time.Second),
		MaxConcurrentRequests:  100,
		RequestVolumeThreshold: 10,
		SleepWindow:            xtime.Duration(5 * time.Second),
		ErrorPercentThreshold:  50,
	}

	g := NewGroup(c)
	if err := g.Do("run", func() error {
		return errors.Wrap(ecode.ServiceUnavailable, "break now")
	}); err != nil {
		if !ecode.EqualError(ecode.ServiceUnavailable, err) {
			t.Error(err)
		}
	}
}

func TestConcurrentDo(t *testing.T) {
	c := &Config{
		Namespace:              "run",
		Timeout:                xtime.Duration(1 * time.Second),
		MaxConcurrentRequests:  200,
		RequestVolumeThreshold: 10,
		SleepWindow:            xtime.Duration(5 * time.Second),
		ErrorPercentThreshold:  50,
	}

	g := NewGroup(c)

	wg := new(sync.WaitGroup)
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := g.Do("concurrent", func() error {
				time.Sleep(time.Millisecond * 500)
				return nil
			}); err != nil {
				t.Logf("hystrix error(%v)", err)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkGroupDo(b *testing.B) {
	c := &Config{
		Namespace:              "run",
		Timeout:                xtime.Duration(1 * time.Second),
		MaxConcurrentRequests:  1000,
		RequestVolumeThreshold: 10,
		SleepWindow:            xtime.Duration(5 * time.Second),
		ErrorPercentThreshold:  50,
	}
	g := NewGroup(c)

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := g.Do("benchmark", func() error {
				time.Sleep(time.Millisecond * 5)
				return errors.Wrap(ecode.ServiceUnavailable, "break now")
			}); err != nil {
				b.Logf("hystrix error(%v)", err)
			}
		}
	})
}
