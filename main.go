package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

var (
	appstate     state
	homeTemplate *template.Template
)

//TemplateData contains data used within the view
type TemplateData struct {
	IsUnhealthy bool
	Hostname    string
	Message     template.HTML
}

//Config for starting app
type Config struct {
	Host         string
	Port         string
	ReadTimout   time.Duration
	WriteTimeout time.Duration
}

type state struct {
	UnhealthyAfter int64
	Requests       int64
}

func (s *state) IsUnhealthy() bool {
	return s.IsUnhealthyMode() && s.Requests > s.UnhealthyAfter
}

func (s *state) IsUnhealthyMode() bool {
	return s.UnhealthyAfter != 0
}

func (s *state) AddRequest() {
	atomic.AddInt64(&s.Requests, 1)
}

// HTMLServer struct
type HTMLServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

func init() {
	homeTemplate = template.Must(template.ParseFiles(path.Join("templates", "page.tpl")))
	flag.Int64Var(&appstate.UnhealthyAfter, "unhealthyAfter", 0, "set the number of request after which the service should fail (returning 500 httpcode)")
}

func main() {
	flag.Parse()

	config := Config{
		Host:         "localhost",
		Port:         "8080",
		ReadTimout:   5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	htmlServer := Start(config)
	defer htmlServer.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	color.New(color.FgYellow).Add(color.Bold).Printf("\nKubia :: Received terminate signal\n")
}

// Start the HTMLServer
func Start(config Config) *HTMLServer {
	http.HandleFunc("/", homeHandler)

	htmlServer := HTMLServer{
		server: &http.Server{
			Addr:           fmt.Sprintf("%s:%s", config.Host, config.Port),
			ReadTimeout:    config.ReadTimout,
			WriteTimeout:   config.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}

	htmlServer.wg.Add(1)

	go func() {
		var unhealthyStr string
		if appstate.IsUnhealthyMode() {
			unhealthyStr = " in Unhealthy Mode"
		}
		color.New(color.FgGreen).Add(color.Bold).Printf("Kubia :: HTMLServer :: Starting%s at 'http://%s:%s'\n", unhealthyStr, config.Host, config.Port)
		htmlServer.server.ListenAndServe()
		htmlServer.wg.Done()
	}()

	return &htmlServer
}

// Stop the HTMLServer
func (htmlServer *HTMLServer) Stop() error {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	color.New(color.FgYellow).Add(color.Bold).Println("Kubia :: HTMLServer :: Server stopping...")

	if err = htmlServer.server.Shutdown(ctx); err != nil {
		if err = htmlServer.server.Close(); err != nil {
			color.New(color.FgRed).Add(color.Bold).Printf("Kubia :: HTMLServer :: Server stopping with error: %v\n", err)
			return err
		}
	}

	htmlServer.wg.Wait()
	color.New(color.FgYellow).Add(color.Bold).Println("Kubia :: HTMLServer: Server stopped")
	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	hostname, err := os.Hostname()
	checkErr(err)

	appstate.AddRequest()

	data := TemplateData{
		IsUnhealthy: appstate.IsUnhealthy(),
		Hostname:    hostname,
		Message:     "",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if data.IsUnhealthy {
		w.WriteHeader(http.StatusInternalServerError)
		data.Message = template.HTML(fmt.Sprintf("<strong>#%d I'm not well. Please restart me!</strong>", appstate.Requests))
		err = homeTemplate.Execute(w, data)
		checkErr(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	data.Message = template.HTML(fmt.Sprintf("#%d You've hit %s", appstate.Requests, hostname))
	err = homeTemplate.Execute(w, data)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
