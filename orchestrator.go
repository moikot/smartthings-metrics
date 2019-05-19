package main

type Orchestrator interface {
	Execute() error
}

func NewOrchestrator(token string) Orchestrator {
	return &orchestrator{
		recorder:  NewMetricRecorder(),
		processor: NewStatusProcessor(),
		reader:    NewStatusReader(token),
	}
}

type orchestrator struct {
	recorder  MetricRecorder
	processor StatusProcessor
	reader    StatusReader
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
		for _, val := range values {
			o.recorder.Record(val)
		}
	}
	return nil
}
