/*
Copyright (c) 2020 Sergey Anisimov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/moikot/smartthings-metrics/recording"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

const (
	defaultPromAddr = ":9153"
)

var (
	exitFunc = os.Exit
	stopChan = make(chan os.Signal, 1)
)

func main() {
	var token, intervalStr string

	flag.StringVar(&token, "token", os.Getenv("API_TOKEN"), "The API token.")
	flag.StringVar(&intervalStr, "interval", os.Getenv("REFRESH_INTERVAL"), "The status refresh interval in seconds.")

	flag.Parse()

	exitFunc(Run(token, intervalStr))
}

func ParseInterval(intervalStr string) (int, error) {
	if len(intervalStr) > 0 {
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			return 0, errors.New("the interval specified is not an integer")
		}
		if interval <= 0 {
			return 0, errors.New("the interval should be greater than zero")
		}
		return interval, nil
	}
	return 60, nil // Default interval is 60 seconds
}

func Run(token, intervalStr string) int {

	if token == "" {
		flag.PrintDefaults()
		return 1
	}

	interval, err := ParseInterval(intervalStr)
	if err != nil {
		_, _ = fmt.Fprint(flag.CommandLine.Output(), err)
		return 1
	}

	go servePrometheus(defaultPromAddr)

	loop := recording.NewLoop(token, interval)
	loop.Start()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan

	return 0
}

func servePrometheus(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("unable to create HTTP server, error: " + err.Error())
	}
}
