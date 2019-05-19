package extractors

import (
	"github.com/SmartThingsOSS/smartapp-go/pkg/smartthings/models"
)

type Extractors map[string]ValueExtractor

func (exts Extractors) Add(e ValueExtractor) {
	exts[e.ID()] = e
}

type Labels map[string]string

type AttributeValue struct {
	Name   string
	Labels Labels
	Value  float64
}

type ValueExtractor interface {
	ID() string
	Name() string
	Labels() []string
	Extract(status models.CapabilityStatus) (*AttributeValue, error)
}
