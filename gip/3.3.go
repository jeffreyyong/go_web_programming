package main

import (
	"compress/gzip"
	"io"
	"os"
)

func main() {
	// Collects a list of files passed in on the command line
	for _, file := range os.Args[1:] {
		compress(file)
	}
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
