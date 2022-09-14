package util

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestThrottle(t *testing.T) {
	k := atomic.Bool{}
	k.Store(true)
	wg := &sync.WaitGroup{}
	calc := atomic.Int32{}
	throttle, _ := GetThrottle(1000*60 + 207)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for k.Load() {
				if throttle.Try() {
					calc.Add(1)
				}
				runtime.Gosched()
			}
			wg.Done()
		}()
	}
	timer := time.NewTimer(4*time.Second + 200*time.Millisecond)
	defer timer.Stop()
	_ = <-timer.C
	k.Store(false)
	wg.Wait()
	processed := calc.Load()
	if processed < 4010 || processed > 4020 {
		t.Error("频率控制精度出现问题 可能与主机性能过低有关")
	}
}
