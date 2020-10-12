package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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

	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})

	interval, err := ParseInterval(intervalStr)
	if err != nil {
		_, _ = fmt.Fprint(flag.CommandLine.Output(), err)
		return 1
	}

	go servePrometheus(defaultPromAddr)

	orchestrator := NewOrchestrator(token, l)

	go func() {
		for {
			err := orchestrator.Execute()
			if err != nil {
				l.Error(err)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()

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
