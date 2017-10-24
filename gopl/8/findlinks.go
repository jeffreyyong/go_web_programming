package main

import (
	"fmt"
	"log"
	"os"

	"go_web_programming/gopl/8/links"
)

// Crawl1 crawls web links starting with the comman-line arguments.
// This version quickly ecahusts available file descriptors due to excessive concurrent calls to links.Extract

// Also, it never terminates because the worklist is never closed

func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

// Resembles a breadthFirst search
// The worklist records the queue of items that need processing, each item being a list of URLs to crawl,
// use channel to represent the queue rather than using a slice

// Each call to crawl occurs in its own goroutine and sends the links it discovers back to the worklist
func main() {
	worklist := make(chan []string)

	// Start with the command-line arguments.
	// Initial send of caommand line argurments to the worklist must run in its own goroutine
	// to avoid deadlock, a stuck situation in which both the main goroutine and a crawler goroutine attempt
	// to send to each other while neither is receiving. Alternative would be to use a buffered channel.
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				// crawl goroutine takes link as an explicit parameter to avoid the problem of loop variable capture
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

// The crawler is highly concurrent and prints a storm of URLs, but it has two problems.
// The initial error message is a DNS lookup failure for a reliable domain.
// The subsequent message message reveals that: the program created too many network connections at once
// that it exceeded the per-process limit on the number of open files, causing operations such as
// DNS lookups and calls to net.Dial to start failing

// Limiting factor of parallelism in a system:
// 1) Number of CPU cores for compute-bound workloads
// 2) Number of spindles and heads for local disk I/O operations
// 3) the bandwidth of the network for streaming downloads
// 4) serving capability of a web service
