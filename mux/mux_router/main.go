package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ArticleHandler(w http.ResponseWriter, r *http.Request) {
	// Mux parses the variables in the path. All the variables that are parsed are available calling
	// the Vars function.
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category is: %v\n", vars["category"])
	fmt.Fprintf(w, "ID is: %v\n", vars["id"])
}

func main() {
	// Create a new router
	r := mux.NewRouter()
	// Attach an elegant path with handler
	r.HandleFunc("/articles/{category}/{id:[0-9]+}", ArticleHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
