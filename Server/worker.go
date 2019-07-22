package main

import (
	_ "expvar"
	"fmt"
	_ "net/http/pprof"
)

// Structure that is sent to worker from a requester
type Job struct {
	Username  string
}


type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	Quit       chan bool
}

// Creating a new worker with and a JobChannel
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		Quit:       make(chan bool)}
}

// Start working --waiting for job to execute
func (w Worker) Start() {
	go func() {
		for {
			// Add my jobQueue to the worker pool.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// Execution of a job
				fmt.Println("Username: ",job.Username)

			case <-w.Quit:
				// Closing connections and quiting, after calling Stop func --bool is true
				close(w.JobChannel)
				close(w.Quit)
				close(w.WorkerPool)
				return
			}
		}
	}()
}

// Setting boolean var as true --quiting
func (w Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}


