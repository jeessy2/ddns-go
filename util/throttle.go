package util

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Throttle interface {
	Try() bool
	Close()
}

type atomicInt64 struct {
	v int64
}

// Load atomically loads and returns the value stored in x.
func (x *atomicInt64) Load() int64 { return atomic.LoadInt64(&x.v) }

// Store atomically stores val into x.
func (x *atomicInt64) Store(val int64) { atomic.StoreInt64(&x.v, val) }

// Add atomically adds delta to x and returns the new value.
func (x *atomicInt64) Add(delta int64) (new int64) { return atomic.AddInt64(&x.v, delta) }

type throttleImp struct {
	times   *atomicInt64
	qpm     int64
	step    int64
	cnt     []int64
	once    *sync.Once
	closeCh chan struct{}
}

func (t *throttleImp) start() {
	t.once.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			second := 0
			for {
				end := false
				select {
				case <-t.closeCh:
					end = true
				case <-ticker.C:
					second += 1
					if second == 60 {
						second = 0
					}
					curStep := t.step + t.cnt[second]
					newVal := t.times.Add(curStep)
					if newVal > t.qpm {
						t.times.Add(-curStep)
					}
				}
				if end {
					break
				}
			}
		}()
	})
}

func (t *throttleImp) Close() {
	t.closeCh <- struct{}{}
}

func (t *throttleImp) Try() bool {
	newVal := t.times.Add(-1)
	if newVal >= 0 {
		return true
	}
	t.times.Add(1)
	return false
}

var errExceedQPM = errors.New("最大只支持2000qps")

func GetThrottle(timesPerMin int64) (Throttle, error) {
	if timesPerMin > 2000*60 {
		return nil, errExceedQPM
	}
	cnt := make([]int64, 60)
	calc, step := int64(0), timesPerMin/60
	for i := 0; i < 60; i++ {
		all := int64(float64(i+1)/60*float64(timesPerMin) - float64(calc))
		cnt[i] = all - step
		calc += all
	}
	throttle := &throttleImp{
		times:   &atomicInt64{},
		step:    step,
		cnt:     cnt,
		once:    &sync.Once{},
		closeCh: make(chan struct{}),
		qpm:     timesPerMin,
	}
	throttle.start()
	return throttle, nil
}
