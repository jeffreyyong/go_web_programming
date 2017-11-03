package main // The du4 variant includes cancellation:
import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// It terminates quickly when the user hits return.

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func main() {
	// Determine the intial directories.
	roots := os.Args[1:]
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Cancel traversal when input is detected.
	go func() {
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		close(done)
	}()

	// Traverse each root of the file tree in parallel.
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

	// Print the results periodically
	tick := time.Tick(500 * time.Millisecond)
	var nfiles, nbytes int64

loop:

	for {
		select {
		case <-done:
			// Drain fileSizes to allow existing goroutines to finish.
			// Discarding all values until the channel is closed.
			// This ensures that any active calls to walkDir can run to completion without getting stuck sending to fileSizes.
			for range fileSizes {
				// Do nothing.
			}
			return
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
}

func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
	defer n.Done()
	// walkDir goroutine polls the cancellation status when it begins, and returns without doing anything if the status is set.
	// This returns all goroutines created after cancellation into no-ops
	if cancelled() {
		return
	}
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subsidr := filepath.Join(dir, entry.Name())
			go walkDir(subdir, n, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

var sema = make(chan struct{}, 20) // concurrency-limiting counting semaphore

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	// The select makes the operation cancellable and reduces the typical cancellation latency of the program from hundres of milliseconds to tens.
	select {
	case sema <- struct{}{}: // acquire token
	case <-done:
		return nil
	}

	defer func() { <-sema }() // release token

	f, err := os.Open(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(0) // 0=> no limit; read all entries

	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)

	}

	return entries
}
