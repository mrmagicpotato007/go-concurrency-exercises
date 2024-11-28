//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"fmt"
	"sync"
	"time"
)

type safeMap struct {
	v  map[string]bool
	mu sync.Mutex
}

var sp *safeMap

func init() {
	sp = &safeMap{v: map[string]bool{}}
}

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup, limiter <-chan time.Time) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	//wg.Add(len(urls))

	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently

		sp.mu.Lock()
		if !sp.v[u] {
			sp.v[u] = true
			sp.mu.Unlock()
			wg.Add(1)
			<-limiter
			go Crawl(u, depth-1, wg, limiter)
		} else {
			sp.mu.Unlock()
		}

	}
	//fmt.Println("parsed all")
}

func main() {
	var wg sync.WaitGroup
	startTime := time.Now()
	wg.Add(1)
	if !sp.v["http://golang.org/"] {
		sp.v["http://golang.org/"] = true
	}
	limiter := time.Tick(time.Second)
	Crawl("http://golang.org/", 4, &wg, limiter)
	wg.Wait()
	fmt.Println(time.Since(startTime))
}
