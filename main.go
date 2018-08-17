package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/fatih/color"
)

const port = ":8080"

var requestCount int32

func main() {
	c := color.New(color.FgGreen).Add(color.Bold)
	_, err := c.Printf("Kubia server starting on port %s...\n", port)
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request from " + r.RemoteAddr)

		hostname, _ := os.Hostname()
		atomic.AddInt32(&requestCount, 1)

		if requestCount > 10 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "I'm not well. Please restart me!")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "#%d You've hit %s\n", requestCount, hostname)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
