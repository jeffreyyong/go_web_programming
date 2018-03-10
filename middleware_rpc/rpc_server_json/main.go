package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	// Derives from the standard net/rpc package but uses a single HTTP request per call instead of
	// persistent connections.
	"github.com/gorilla/rpc"
	rpcjson "github.com/gorilla/rpc/json"
)

// Create the Args and Books struct to hold the information about the JSON arguments passed and the book structure
type Args struct {
	Id string
}

type Book struct {
	Id     string `"json:string,omitempty"`
	Name   string `"json:name,omitempty"`
	Author string `"json:author,omitempty"`
}

// JSONServer has the GiveBookDetail remote function. This struct is a service created to register with the
// `Register Service` function of the RPC server. Also registering thed codec as JSON.
type JSONServer struct{}

// Receives the request, it reads the file from the filesystem and parses it
// If the given ID matches any book, then the server sends the information back to the client in the JSON format
// The reply reference is passed to the remote function. In the remote function, setting the value of the reply
// with the matched book. If the ID sent by the client matches with any of the books in the JSON< then the data is filled
// If there is no match, then empty data will be sent back by the RPC server.
func (t *JSONServer) GiveBookDetail(r *http.Request, args *Args, reply *Book) error {
	var books []Book
	raw, readerr := ioutil.ReadFile("./books.json")
	if readerr != nil {
		log.Println("error:", readerr)
		os.Exit(1)
	}

	marshalerr := json.Unmarshal(raw, &books)
	if marshalerr != nil {
		log.Println("error:", marshalerr)
		os.Exit(1)
	}

	// Iterate over JSON data to find the given book
	for _, book := range books {
		if book.Id == args.Id {
			*reply = book
			break
		}
	}
	return nil
}

func main() {
	s := rpc.NewServer()
	// Codec is chosen based on the Content-Type header from the request.
	s.RegisterCodec(rpcjson.NewCodec(), "application/json")
	// Service methods also receive the http.Request as a parameter.
	s.RegisterService(new(JSONServer), "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	http.ListenAndServe(":1234", r)
}
