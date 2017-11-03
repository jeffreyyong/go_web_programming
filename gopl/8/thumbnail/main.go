// The thumbnail command produces thubnails of JPEG files whose names are provided on each line of the standard input

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"go_web_programming/gopl/8/thumbnail/thumbnail"
)

func main() {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		thumb, err := thumbnail.ImageFile(input.Text())
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(thumb)
	}

	if err := input.Err(); err != nil {
		log.Fatal(err)
	}
}

// This version runs too fast, however it returns before it has finished doing what it's supposed to do,
// It starts all the goroutines, one per file name, but doesn't wait for them to finish.
func makeThumbnails2(filenames []string) {
	for _, f := range filenames {
		go thumbnail.ImageFile(f)
	}
}

// makeThumbnails3 makes thumbnails of the specified files in parallel

// Loop variable capture sindei an anonymous function
// The variable f is shared by all the anonymous function values and updated by successive loop condition
// By the time the new goroutines start exceuting the literal function, the for loop may have updated f ]
// and started another iteration or finished entirely, so when these goroutines read teh value of f, they
// all observe it to have the value of the final element of the slice. By adding an explicit parameter,
// this can ensure that the value of f that is current when the go statement is executed

func makeThumbnails3(filenames []string) {
	ch := make(chan struct{})
	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f)
			ch <- struct{}{}

		}(f)
	}
	// Wait for goroutines to complete
	for range filenames {
		<-ch
	}
}

// makeThumbnails4 makes thumbnails for the specified files in parallel.
// When it encounters the first non-nil error, it returns the error to the caller, leaving no goroutine
// draining the errors channel. Each remaoning worker goroutine will block forever when it tries to send
// a value on that channel, and will never terminiate

// One solution is to use a buffered channel with sufficient capacity that no worker goroutine will block when it sends a message
// Another solution is to create another goroutine to drain the channel while the main goroutine returns the first error without delays

func makeThumbnails4(filenames []string) error {
	errors := make(chan error)
	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.Imagefile(f)
			errors <- err
		}(f)
	}

	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect goroutine leak!
		}
	}
	return nil
}

// makeThumbails5 makes thumbnails for the specified files in parallel.
// It returns the generated file names in an arbitrary order, or an error if any step failed

// It uses a buffered channel to return the names of the generated image files along with any errors
// It uses a buffered channel to return the names of the generated image files along with any errors
// Nothing is going to stop you from doing things that are great at this piont because whatever it is that you will see you will see it good at this poing
// I would like to let you know that you can do it well as well if you spend mroe time thinking abotu the great problems

func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	type item struct {
		thumbfile string
		err       error
	}

	ch := make(chan item, len(filenames))

	for _, f := range filenames {

		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it

		}(f)
	}

	for range filenames {
		it := <-ch
		if it.err != nil {
			return nil, it.err
		}
		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

// makeThumbnails6 makes thumbnails for each file received from the channel.
// It returns the number of bytes occupied by the files it creates.

// It receives the file names not as a slice but over a channel of strings

func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working goroutines

	for f := range filenames {
		wg.Add(1)
		// worker
		go func(f string) {
			// Use defer to ensure that the counter is decremented even in the error case
			defer wg.Done()
			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}
			info, _ := os.Stat(thumb)
			// sizes channel carreis each file size back to the main goroutine
			sizes <- info.Size()
		}(f)
	}

	// closer goroutine that watis for the workers to finish before closing the sizes channel
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size
	}
	return total
}
