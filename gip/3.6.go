package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	w := newWords()
	for _, f := range os.Args[1:] {
		wg.Add(1)
		go func(file string) {
			if err := tallyWords(file, w); err != nil {
				fmt.Println(err.Error())
			}
			wg.Done()
		}(f)
	}
	wg.Wait()

	fmt.Println("Words that appear more than once:")
	// Locks and unlocks the map whe iterate at the end
	// Strictly speaking, this isn't necessary because this section won't
	// happen until all files are processed.
	w.Lock()
	for word, count := range w.found {
		if count > 1 {
			fmt.Printf("%s: %d\n", word, count)
		}
	}
	w.Unlock()
}

type words struct {
	sync.Mutex
	found map[string]int
}

func newWords() *words {
	return &words{found: map[string]int{}}
}

func (w *words) add(word string, n int) {
	// Locks the object, modifies the map, and then unlocks the object
	w.Lock()
	defer w.Unlock()
	count, ok := w.found[word]
	if !ok {
		w.found[word] = n
		return
	}
	w.found[word] = count + n
}

func tallyWords(filename string, dict *words) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scnaner.Scan() {
		word := strings.ToLower(scanner.Text())
		dict.add(word, 1)
	}

	return scanner.Err()
}

// Inside the add method, you lock the object, modify the map, and then unlock the object
// When multiple goroutines enter the ad dmethod, the first will get the lock
// and others will wait until the lock is released. This will prevent multiple goroutines
// from modifying the map at the same time

// It's important to note that locks work only when all access to the data is managed by the same lock.
// If some data is accessed with locks, and others without, a race condition can still occur.
