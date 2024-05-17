package api

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const N = 50000

/*
=== RUN   TestBitkv_Put
    api_test.go:34: put use 2.223114041s
--- PASS: TestBitkv_Put (8.71s)
=== RUN   TestBitkv_Get
    api_test.go:44: get use 2.203679042s
--- PASS: TestBitkv_Get (5.44s)
PASS
*/

var c = NewClient("127.0.0.1:7070")

func TestBitkv_Put(t *testing.T) {
	logx.Disable()
	keys := randStrings()
	values := randStrings()
	now := time.Now()
	for i := range N {
		c.Put(keys[i], values[i])
	}
	t.Log("put use", time.Since(now))
}

func TestBitkv_Get(t *testing.T) {
	logx.Disable()
	keys := randStrings()
	now := time.Now()
	for i := range N {
		c.Get(keys[i])
	}
	t.Log("get use", time.Since(now))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStrings() (strings []string) {
	var mu sync.Mutex
	var wg sync.WaitGroup

	wg.Add(N)

	for range N {
		go func() {
			defer wg.Done()
			b := make([]byte, 25)
			for i := range b {
				b[i] = letterBytes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterBytes))]
			}
			mu.Lock()
			strings = append(strings, string(b))
			mu.Unlock()
		}()
	}

	wg.Wait()
	return
}
