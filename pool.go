package worker

import (
	"sync/atomic"
)

type Job func()

type Pool struct {
	jobs chan workRequest

	workers []*workerWrapper

	queuedJobs int64
}

func NewPool(n int) *Pool {
	p := &Pool{
		jobs: make(chan workRequest),
	}
	var workers = make([]*workerWrapper, n)
	for i := 0; i < n; i++ {
		workers[i] = newWorkerWrapper(p.jobs)
	}
	p.workers = workers
	return p
}

func (p *Pool) Send(j Job) {

	atomic.AddInt64(&p.queuedJobs, 1)

	jch := <-p.jobs

	jch.jobChan <- j

	atomic.AddInt64(&p.queuedJobs, -1)

}
