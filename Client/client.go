package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	URL = "http://localhost:8080/username"
	contentType = "application/json"
)

var (
	i int32 = 1
	reqNum, errReq = strconv.Atoi(os.Getenv("REQ_PER_MIN")) // Environment Variables
)

type Response struct {
	StatusCode int
}

type Json struct {
	Username string
}

func ExecuteRequests() {
	var wg sync.WaitGroup // Goroutines synchronization
	procNum := runtime.NumCPU()

	runtime.GOMAXPROCS(procNum) // Set GOMAXPROCS in case of using GO version lower than 1.5 -- default value is 1
	// --in some cases it is better to use value 1

	wg.Add(procNum)
	progTime := time.Now()

	// Goroutines division on cores -- Parallel computing
	for p := 0; p < procNum; p++ {
		go func() {
			defer wg.Done()

			limiter := time.Tick((time.Minute /time.Duration(reqNum/procNum))) // RateLimit

			for time.Since(progTime).Seconds()<=60 {
				<-limiter
				go PostRequest()
			}
		}()
	}
	wg.Wait() // Blocking till the end of goroutines work
}

func PostRequest() {
	response := make(chan Response)
	buf, err := json.Marshal(Json{Username: String(6)}) // JSON coding
	if err != nil{
		println("Error during JSON coding")
	}

	go func() { // Waiting for response
		res := <-response
		if(res.StatusCode == 201){
			atomic.AddInt32(&i,1) // Increment counter
		}else{
			log.Println(res.StatusCode)
		}
	}()

	resp, err := http.Post(URL, contentType, bytes.NewBuffer(buf)) // Post request
	if err != nil {
		log.Fatalln(err)
		return
	}

	response <- Response{resp.StatusCode}
	defer resp.Body.Close()
}


// Generating a random string for a JSON structure
func String(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func main() {
	if errReq != nil {
		log.Fatalln("Enter right Environment Variables") // Handle convert error
	}

	ExecuteRequests()
	println(i) // Print
}