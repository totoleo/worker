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
	runtime.GOMAXPROCS(10)
	go func() {

		pool := NewPool(50)

		time.Sleep(1 * time.Second)
		//w := &sync.WaitGroup{}
		//w.Add(100)

		for i := 0; i < 100; i++ {
			var p func()
			index := i
			if i%2 == 0 {
				p = func() {
					defer func() {
						err := recover()
						if err != nil {
							fmt.Println(err)
						}
					}()
					fmt.Println(index, "hello:"+strconv.Itoa(index))
					panic("err")
				}
			} else {
				p = func() {
					time.Sleep(1 * time.Second)
					fmt.Println(index, time.Now())
				}
			}
			pool.Send(p)
		}
	}()

	ticker := time.NewTicker(200 * time.Millisecond)
	for true {
		select {
		case <-ticker.C:
			fmt.Println("num of routine:", runtime.NumGoroutine())
		}
	}
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
