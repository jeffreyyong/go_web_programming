// Outyet is a web server that announces whether or not a particular Go version has been tagged.
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	httpAddr   = flag.String("http", "localhost:8080", "Listen address")
	poolPeriod = flag.Duration("poll", 5*time.Second, "Poll period")
	version    = flag.String("version", "1.9", "Go version")
)

const baseChangeURL = "https://code.google.com/p/go/source/detail?r="

func main() {
	flag.Parse()
	changeURL := fmt.Sprintf("%sgo%s", baseChangeURL, *version)
	http.Handle("/", NewServer(*version, changeURL, *pollPeriod))
}
