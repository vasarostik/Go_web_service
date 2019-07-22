package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	MaxQueueTest = 100
	JobQueueTest = make(chan Job, MaxQueueTest)
)



// 1. Check  exact duration of sending 1 000 000 requests:
// go test -bench=BenchmarkRequest -benchtime=60s

func SendRequest(URL string) {
	buf, err := json.Marshal(Job{Username: "test"})
	if err != nil{
		panic(err)
	}

	res, err := http.Post(URL, "application/json", bytes.NewBuffer(buf) )
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 201 {
		panic("Inappropriate Status Code")
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = res.Body.Close()
	if err != nil {
		panic(err)
	}
}

func BenchmarkRequest(b *testing.B) {

	for n := 0; n < b.N; n++ {
		SendRequest( "http://127.0.0.1:8080/username")
	}
}

// 2. Check performance, using different number of workers

func BenchmarkHandlingWithOneWorker(b *testing.B) {
	dispatcher := NewDispatcher(JobQueueTest, 1)
	dispatcher.Run()

	for n := 0; n < b.N; n++ {
		JobQueueTest <- Job{Username: "test"}
		println(len(JobQueueTest))
	}
}

func BenchmarkHandlingWithEightWorkers(b *testing.B) {
	dispatcher := NewDispatcher(JobQueueTest, 8)
	dispatcher.Run()

	for n := 0; n < b.N; n++ {
		JobQueueTest <- Job{Username: "test"}
		println(len(JobQueueTest))
	}
}

// 3. Check how many jobs can be performed
// go test -bench=BenchmarkJobQueue -benchtime=3s
func BenchmarkJobQueue(b *testing.B) {
	dispatcher := NewDispatcher(JobQueueTest, 8)
	dispatcher.Run()

	for n := 0; n < b.N; n++ {
		JobQueueTest <- Job{Username: "test"}
		println(len(JobQueueTest))
	}
	close(dispatcher.JobQueue)
	close(dispatcher.WorkerPool)
}





