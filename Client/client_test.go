package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Output struct {
	Username string
}

func benchmarkCkeckProcesses(i int, b *testing.B) {
	var wg sync.WaitGroup

	// set GOMAXPROCS in case of using GO version lower than 1.5 -- default value is 1
	runtime.GOMAXPROCS(i)

	wg.Add(i)

	for p := 0; p < i; p++ {
		go func() {
			defer wg.Done()

			// RateLimit
			limiter := time.Tick((time.Minute / time.Duration(1000000/i)))

			for n := 0; n < b.N; n++ {
				<-limiter
				go PostRequest()
			}
		}()
	}
	wg.Wait()
}

func BenchmarkJSONMarshal(b *testing.B) {
	obj := Json{Username: String(6)}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkJSONUnmarshal(b *testing.B) {
	out, err := json.Marshal(Json{Username: String(6)})
	if err != nil {
		panic(err)
	}

	obj := &Output{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err = json.Unmarshal(out, obj)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkRequestWithOneProcess(b *testing.B) { benchmarkCkeckProcesses(1, b) }

// after BenchmarkRequestWithOneProcess(8,b) we should  multiply the result by 8,
// because it is parallel computing.
func BenchmarkRequestWithEightProcesses(b *testing.B) { benchmarkCkeckProcesses(8, b) }

//Unit Test - go test
func TestPostRequest(t *testing.T) {
	response := make(chan Response)
	buf, _ := json.Marshal(Json{Username: String(6)})            // JSON coding
	resp, _ := http.Post(URL, contentType, bytes.NewBuffer(buf)) // Post request
	assert.NotNil(t, resp, "Response shouldn`t be nil")
	go func() { // Waiting for response
		res := <-response
		if assert.Equal(t, 201, resp.StatusCode, "Status code 201 is expected") {
			atomic.AddInt32(&i, 1) // Increment counter
		} else {
			log.Println(res.StatusCode)
		}
	}()

	response <- Response{resp.StatusCode}
	defer resp.Body.Close()
}

func TestString(t *testing.T) {
	body := String(5)
	assert.Equal(t, 5, len(body), "Body should have length 5")
	assert.NotNil(t, body, "Body shouldn`t be nil")
}
