package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

var (
	MaxWorker, errWorker = strconv.Atoi(os.Getenv("MAX_WORKERS")) // Environment Variables
	MaxQueue,  errQueue  = strconv.Atoi(os.Getenv("MAX_QUEUE"))
	JobQueue = make(chan Job, MaxQueue) // Channel for passing jobs to Dispatcher
)

func RequestHandler(w http.ResponseWriter, r *http.Request, JobQueue chan Job) {
	if r.Method == http.MethodPost {
		job := Job{} // Create Job and push into the jobQueue

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Couldn't parse body", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &job) // JSON decoding
		if err != nil {
			log.Fatal("Decoding error: ", err)
		}

		JobQueue <- job // Passing job to Dispatcher for being executed by a worker

		w.WriteHeader(http.StatusCreated)
		return
	}else{
		w.Header().Set("Allow", "POST")
		http.Error(w, "POST method only", http.StatusMethodNotAllowed)
	}
}


func main() {
	if (errWorker != nil) || (errQueue != nil) {
		log.Fatalln("Enter right Environment Variables")
	}
	// Start the dispatcher
	dispatcher := NewDispatcher(JobQueue, MaxWorker)
	dispatcher.Run()

	// Start the HTTP handler
	http.HandleFunc("/username", func(w http.ResponseWriter, r *http.Request) {
		RequestHandler(w, r, JobQueue)})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
