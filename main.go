package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	exitFunc = os.Exit
	stopChan = make(chan os.Signal, 1)
)

func main() {
	var token, intervalStr string

	flag.StringVar(&token, "token", os.Getenv("API_TOKEN"), "The API token.")
	flag.StringVar(&intervalStr, "interval", os.Getenv("REFRESH_INTERVAL"), "The status refresh interval in seconds.")

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

	orchestrator := NewOrchestrator(token)

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
