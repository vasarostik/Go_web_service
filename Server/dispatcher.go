package main
import 	_ "net/http/pprof"


// Dispatcher
type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	MaxWorkers int
}

// Creating a new dispatcher object and a worker pool
func NewDispatcher(JobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		JobQueue:   JobQueue,
		MaxWorkers: maxWorkers,
		WorkerPool: workerPool,
	}
}

// Initialization and start working of new workers in worker pool,
// Calling dispatch() as a separate hanging goroutine
func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

// Waiting for new jobs in JobQueue and passing them to workers by JobChannel
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.JobQueue:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}
