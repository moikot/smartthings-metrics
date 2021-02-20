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

package recording

import (
	"time"

	"github.com/moikot/smartthings-metrics/extracting"
	"github.com/moikot/smartthings-metrics/readers"
	log "github.com/sirupsen/logrus"
)

func NewLoop(token string, interval uint) *Loop {

	return &Loop{
		DeviceReader:    readers.NewDeviceReader(token),
		StatusProcessor: extracting.NewStatusProcessor(),
		MetricRecorder:  NewMetricRecorder(),
		interval:        interval,
	}
}

type Loop struct {
	readers.DeviceReader
	extracting.StatusProcessor
	MetricRecorder
	interval uint
}

func (l *Loop) Start() {
	go func() {
		for {
			err := l.record()
			if err != nil {
				log.Error(err)
			}
			time.Sleep(time.Duration(l.interval) * time.Second)
		}
	}()
}

func (l *Loop) record() error {
	statuses, err := l.ReadStatuses()
	if err != nil {
		return err
	}

	l.Record(l.GetMeasurements(statuses))
	return nil
}
