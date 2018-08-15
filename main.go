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
	c.Printf("Kubia server starting on port %s...\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request from " + r.RemoteAddr)
		w.WriteHeader(http.StatusOK)

		hostname, _ := os.Hostname()

		fmt.Fprintf(w, "You've hit %s\n", hostname)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
