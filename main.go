package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/fatih/color"
)

const (
	port = ":8080"
	tpl  = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <link rel="icon" href="data:,">
    <title>Kubia</title>
</head>
<body>
	<span align="center">
		<h1>Kubia</h1>
		<p>{{ .Message }}</p>
	</span>
</body>
</html>
`
)

var (
	requestCount   int64
	unhealthyAfter int64
)

type templateData struct {
	Unhealthy bool
	Hostname  string
	Message   string
}

func init() {
	flag.Int64Var(&unhealthyAfter, "unhealthyAfter", 0, "set the number of request after which the service should fail (returning 500 httpcode)")
}

func main() {
	flag.Parse()

	var unhealthyStr string
	if unhealthyAfter != 0 {
		unhealthyStr = " in unhealthy mode"
	}

	_, err := color.New(color.FgGreen).Add(color.Bold).Printf("Kubia server starting%s on port 'http://localhost%s'...\n", unhealthyStr, port)
	checkErr(err)

	t := template.Must(template.New("page").Parse(tpl))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		hostname, err := os.Hostname()
		checkErr(err)

		atomic.AddInt64(&requestCount, 1)

		data := templateData{
			Unhealthy: unhealthyAfter != 0 && requestCount > unhealthyAfter,
			Hostname:  hostname,
			Message:   "",
		}

		if data.Unhealthy {
			w.WriteHeader(http.StatusInternalServerError)

			data.Message = fmt.Sprintf("#%d I'm not well. Please restart me!", requestCount)

			err = t.Execute(w, data)
			checkErr(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		data.Message = fmt.Sprintf("#%d You've hit %s", requestCount, hostname)
		err = t.Execute(w, data)
		checkErr(err)
	})

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
