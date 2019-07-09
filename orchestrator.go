package main

// Orchestrator reads device statuses, processes them converting
// to records containing a name, a value and a set of labels.
// Finally, it passes the resulting records to the recorder.
type Orchestrator interface {
	Execute() error
}

func NewOrchestrator(token string) Orchestrator {
	return &orchestrator{
		recorder:  NewMetricRecorder(),
		processor: NewStatusProcessor(),
		reader:    NewDeviceReader(token),
	}
}

type orchestrator struct {
	recorder  MetricRecorder
	processor StatusProcessor
	reader    DeviceReader
}

func (o *orchestrator) Execute() error {
	statuses, err := o.reader.Read()
	if err != nil {
		return err
	}

	for _, status := range statuses {
		values, err := o.processor.Process(status.Device, status.Status)
		if err != nil {
			return err
		}
		for _, value := range values {
			o.recorder.Record(value)
		}
	}
	return nil
}
