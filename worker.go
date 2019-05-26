package worker

//------------------------------------------------------------------------------

// workRequest is a struct containing context representing a workers intention
// to receive a work payload.
type workRequest struct {
	// jobChan is used to send the payload to this worker.
	jobChan chan<- Job
}

//------------------------------------------------------------------------------

// workerWrapper takes a Worker implementation and wraps it within a goroutine
// and channel arrangement. The workerWrapper is responsible for managing the
// lifetime of both the Worker and the goroutine.
type workerWrapper struct {
	// reqChan is NOT owned by this type, it is used to send requests for work.
	reqChan chan<- workRequest
}

func newWorkerWrapper(
	reqChan chan<- workRequest,
) *workerWrapper {
	w := workerWrapper{
		reqChan: reqChan,
	}

	go w.run()

	return &w
}

func (w *workerWrapper) run() {
	jobChan := make(chan Job)

	for {
		select {
		case w.reqChan <- workRequest{
			jobChan: jobChan,
		}:
			select {
			case payload := <-jobChan:
				payload()
			}
		}
	}
}

//------------------------------------------------------------------------------
