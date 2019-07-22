package main

import (
	"encoding/json"
	"runtime"
	"sync"
	"testing"
	"time"
)

type Output struct{
	Username string
}

func benchmarkCkeckProcesses(i int,b *testing.B) {
	var wg sync.WaitGroup

	// set GOMAXPROCS in case of using GO version lower than 1.5 -- default value is 1
	runtime.GOMAXPROCS(i)

	wg.Add(i)

	for p := 0; p < i; p++ {
		go func() {
			defer wg.Done()

			// RateLimit
			limiter := time.Tick((time.Minute /time.Duration(1000000/i)))

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

func BenchmarkRequestWithOneProcess(b *testing.B) {benchmarkCkeckProcesses(1,b)}
// after BenchmarkRequestWithOneProcess(8,b) we should  multiply the result by 8,
// because it is parallel computing.
func BenchmarkRequestWithEightProcesses(b *testing.B) {benchmarkCkeckProcesses(8,b)}
