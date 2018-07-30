package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/fatih/color"
)

const port = ":8080"

// TODO TEST using os/user and ask for current user

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	c := color.New(color.FgGreen).Add(color.Bold)
	c.Printf("Kubia server starting on port %s as user %s...\n", port, user.Username)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request from " + r.RemoteAddr)
		w.WriteHeader(http.StatusOK)

		hostname, _ := os.Hostname()

		fmt.Fprintf(w, "You've hit %s\n", hostname)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
