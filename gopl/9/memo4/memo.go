// Package memo provides a concurrency-safe memoization a function of a function.
// Requests for different keys proveed in parallel.
// Concurrent request for the same key block until the first completse.
// This implementation uses a Mutex.

package memo

import "sync"

// Func is the type of the function to memoize.
type Func func(string) (interface{}, error)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition/

		e = &entry{ready: make(chan struct{})}

		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		// The closing of the reayd channel happens before any other goroutine receives the broadcast even,
		// so the write to those variabels in the first goroutine happens before they are read by subsequent goroutines
		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this ky.
		memo.mu.Unlock()
		<-e.ready // Wait for ready condition
	}

	return e.res.value, e.res.err
}
