package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	// A waitGroup doesn't need to be initialized.
	var wg sync.WaitGroup
	var i int = -1
	var file string
	for i, file = range os.Args[1:] {
		// For every file added, tell the wait group that it's waiting for one more compress operation
		wg.Add(1)
		// This function calls cmopress and then notifies the wait group that it's done
		go func(filename string) {
			compress(filename)
			wg.Done()
		}(file) // Calling a goroutine in a for loop, need to do this for parameter passing
	}
	wg.Wait() // This waits until all the compressing goroutines have called wg.Done
	fmt.Printf("Compressed %d files\n", i+1)
}

func compress(filename string) error {
	// Open the source file for reading
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	// Open the destination file with the .gz extension added to the source file's name
	out, err := os.Create(filename + ".gz")
	if err != nil {
		return err

		defer out.Close()
	}

	// The gzip.Writer compresses data and then writes it to the underlying file
	gzout := gzip.NewWriter(out)
	// The io.Copy funciton does all the copying
	_, err = io.Copy(gzout, in)
	gzout.Close()

	return err
}

// Instead of closing over file inside the clouser, the file is passed into the programme as filename.
// The variable file is scoped to the for loop, which means that its value will change on each iteration of the loop
// But declaring a goroutine doesn't result in immediate execution
// The the loop runs five times, five goroutines scheduled, but possibily none of them executed.

// On each of those five iterations, the value of file will change.
// By the time the goroutines execute, they may all have the same (fifth) version of the file string
// For file to be scheduled with that iteration's value of file, pass it as a function parameter, which ensures that
// the value of file is passed to each goroutine as it's scheduled
