package main

import "github.com/sirupsen/logrus"

// Orchestrator reads device statuses, processes them converting
// to records containing a name, a value and a set of labels.
// Finally, it passes the resulting records to the recorder.
type Orchestrator interface {
	Execute() error
}

func NewOrchestrator(token string, log logrus.FieldLogger) Orchestrator {
	return &orchestrator{
		recorder:  NewMetricRecorder(log),
		processor: NewStatusProcessor(log),
		reader:    NewDeviceReader(token),
	}
}

type orchestrator struct {
	recorder  MetricRecorder
	processor StatusProcessor
	reader    DeviceReader
	log logrus.FieldLogger
}

func (o *orchestrator) Execute() error {
	statuses, err := o.reader.Read()
	if err != nil {
		return err
	}

	for _, status := range statuses {
		values := o.processor.Process(status.Device, status.Status)
		for _, value := range values {
			o.recorder.Record(value)
		}
	}
	return nil
}
