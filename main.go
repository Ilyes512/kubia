package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/fatih/color"
)

const port = ":8080"

var (
	requestCount int32
	unhealthy    bool
)

func main() {
	flag.BoolVar(&unhealthy, "unhealthy", false, "set this flag to start kubia in unhealthy mode (after 10 requests it returns 500 header")
	flag.Parse()

	var unhealthyStr string
	if unhealthy {
		unhealthyStr = " in unhealthy mode"
	}

	_, err := color.New(color.FgGreen).Add(color.Bold).Printf("Kubia server starting%s on port %s...\n", unhealthyStr, port)
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request from " + r.RemoteAddr)

		hostname, err := os.Hostname()
		if err != nil {
			log.Panic(err)
		}
		atomic.AddInt32(&requestCount, 1)

		if unhealthy && requestCount > 10 {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = fmt.Fprintf(w, "I'm not well. Please restart me!")
			if err != nil {
				log.Panic(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "#%d You've hit %s\n", requestCount, hostname)
		if err != nil {
			log.Panic(err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
