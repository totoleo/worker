package worker

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	pool := NewPool(50)

	go func() {
		t := time.NewTicker(100 * time.Millisecond)
		for true {
			select {
			case <-t.C:
				fmt.Println(runtime.NumGoroutine())
			}
		}
	}()

	w := &sync.WaitGroup{}
	w.Add(100)

	for i := 0; i < 100; i++ {
		var p func()
		index := i
		if i%2 == 0 {
			p = func() {
				defer w.Done()
				fmt.Println(index, "hello:"+strconv.Itoa(index))
			}
		} else {
			p = func() {
				defer w.Done()
				time.Sleep(1 * time.Second)
				fmt.Println(index, time.Now())
			}
		}
		pool.Send(p)
	}

	w.Wait()

}

type job struct {
	data interface{}
	w    *sync.WaitGroup
}

func (j *job) On() {

	switch j.data.(type) {
	case string:
	case time.Time:
		time.Sleep(1 * time.Second)
	}
	j.w.Done()
}
