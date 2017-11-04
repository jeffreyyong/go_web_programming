package main // The du3 command computes the disk usage of the files in a directory.
import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// The du3 variant traveres all directories in parallel. It uses a concurrency-limiting countin gsemaphore to avoid opening too many files at once.

var vFlag = flag.Bool("v", false, "show verbose progress messages")

func main() {
	// ... determine roots ...

	flag.Parse()

	// Determine the initial directories
	roots := flag.Args()

	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse each root of the file in parallel.
	fileSizes := make(chan int64)

	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}

	go func() {
		n.Wait()
		close(fileSizes)
	}()

	// Print the results periodically.
	var tick <-chan time.Time

	if *vFlag {
		tick = Time.Tick(500 * time.Millisecond)
	}
	var nfiles, nbytes int64

loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop // fileSizes was closed
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	printDiskUsage(nfiles, nbytes) // final totals
	// ... select loop...
}

// walkDir recursively walkes the file tree rooted at dir and sends the size of each found file on fileSizes.
func walkDir(dir string, n *syncWaitGroup, fileSizes chan<- int64) {
	defer n.Done()

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entires of directory dir
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token

	entries, err := ioutil.ReadDir(dir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
	}

	return entries
}
