// This version uses a buffered channel as a counting semaphore to limit the number of concurrent calls to links.Extract

package main

import (
	"fmt"
	"go_web_programming/gopl/8/links"
	"log"
	"os"
)

// tokens is a counting sempahore used to enforce a limit of 20 concurrent requests.
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
	fmt.Println(url)

	// Good practice to keep the semaphore operation as close as possible to the I/O operation it regulates
	tokens <- struct{}{} //acquire a token
	list, err := links.Extract(url)
	<-tokens //release the token

	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	worklist := make(chan []string)
	// keeps track of the number of sends to the worlist that are yet to occur
	var n int // number of pending sends to worklist

	// start with the command-line arguments.
	n++
	go func() { worklist <- os.Args[1:] }()

	// Crawl the web concurrently.
	seen := make(map[string]bool)
	// The main loop terminates when n falls to zero, since there is no more work to be done
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				n++
				go func(link string) {
					worklist <- crawl(link)
				}(link)
			}
		}
	}

}

// The second problem si that the program never terminates, even when it has discovered all the links reachable form the initial URLs
