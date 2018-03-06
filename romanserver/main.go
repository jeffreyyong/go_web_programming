package main

import (
	"fmt"
	"go_web_programming/romanserver/roman_numerals"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPathElements := strings.Split(r.URL.Path, "/")

		if urlPathElements[1] == "roman_number" {
			number, _ := strconv.Atoi(strings.TrimSpace(urlPathElements[2]))
			if number == 0 || number > 10 {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - Not Found\n"))
			} else {
				fmt.Fprintf(w, "%q\n", html.EscapeString(roman_numerals.Numerals[number]))
			}
		} else {
			// For all other requests, tell that client sent a bad request
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Bad requests\n"))
		}
	})

	// Create a server and run it on port 8000
	s := &http.Server{
		Addr:           ":8000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Listening...")
	s.ListenAndServe()
}
