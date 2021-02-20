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
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/moikot/smartthings-metrics/recording"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultPromAddr = ":9153"
)

var (
	exitFunc = os.Exit
	stopChan = make(chan os.Signal, 1)
)

func main() {

	token := flag.String("token", "", "The API token.")
	interval := flag.Uint("interval", 60, "The status refresh interval in seconds.")
	verbosity := flag.String("verbosity", "info", "The verbosity level (error, warn, info, debug, trace).")

	flag.Parse()

	args, err := validateArgs(*token, *interval, *verbosity)
	if err != nil {
		fmt.Fprint(flag.CommandLine.Output(), "error: ", err, "\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: smartthings-metrics [flags]:\n")
		flag.PrintDefaults()
		exitFunc(1)
	}

	Run(*args)
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (log.Level, error) {
	switch strings.ToLower(lvl) {
	case "error":
		return log.ErrorLevel, nil
	case "warn":
		return log.WarnLevel, nil
	case "info":
		return log.InfoLevel, nil
	case "debug":
		return log.DebugLevel, nil
	case "trace":
		return log.TraceLevel, nil
	}

	var l log.Level
	return l, fmt.Errorf("not a valid verbosity level: %q", lvl)
}

func Run(args args) {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(args.level)

	go servePrometheus(defaultPromAddr)

	loop := recording.NewLoop(args.token, args.interval)
	loop.Start()

	signal.Notify(stopChan, os.Interrupt, os.Kill)
	<-stopChan
}

type args struct {
	token    string
	interval uint
	level    log.Level
}

func validateArgs(token string, interval uint, verbosity string) (*args, error) {
	if token == "" {
		return nil, errors.New("personal access token is not defined")
	}

	level, err := ParseLevel(verbosity)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse verbosity level")
	}

	return &args{
		token:    token,
		interval: interval,
		level:    level,
	}, nil
}

func servePrometheus(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("server started at port %v", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic("unable to create HTTP server, error: " + err.Error())
	}
}
