package main

import (
	"fmt"
	"net/http"
)

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middleware before request phase!")
		// Pass control back to the handler, it handles the handler logic which is the mainLogic
		handler.ServeHTTP(w, r)
		fmt.Println("Executing middleware after response phase!")
	})
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	// Business logic goes here
	fmt.Println("Executing mainHandler...")
	w.Write([]byte("OK"))
}

func main() {
	// HandlerFunc returns a HTTP Handler
	mainLogicHandler := http.HandlerFunc(mainLogic)
	http.Handle("/", middleware(mainLogicHandler))
	http.ListenAndServe(":8000", nil)
}

// The http.Handle function expects an HTTP handler. Wrap up the logic in such a way that, a handler gets returned
// but the execution is modified

// Passing the main handler into the middleware. Then middleware takes it and returns a function embedding this
// main handler logic in it. This makes all the request coming to the handler pass through the middleware logic.
