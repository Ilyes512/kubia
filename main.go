package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
)

const port = ":8080"

func main() {
	c := color.New(color.FgGreen).Add(color.Bold)
	_, err := c.Printf("Kubia server starting on port %s...\n", port)
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request from " + r.RemoteAddr)

		hostname, _ := os.Hostname()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "You've hit %s\n", hostname)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
