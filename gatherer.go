package funnel

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type Gatherer interface {
	Gather(DataPoint) error
	Accepts(StreamType) bool
}

// NewDataDogBackend returns a Gatherer which collects datapoints
// to the specified DataDog address. It uses the datadog statsd package
// with the buffered client.
func NewDataDogBackend(addr string, prefix string) (Gatherer, error) {
	var c *statsd.Client
	var err error

	c, err = statsd.NewBuffered(addr, 128)
	if err != nil {
		return nil, err
	}

	c.Namespace = prefix + "."
	ddb := &dataDogBackend{
		c,
		make(chan DataPoint),
		GaugeType | CountType | EventType | HistogramType | SetType | SimpleEventType,
	}

	return ddb, nil
}

type dataDogBackend struct {
	client       *statsd.Client
	dataPointsIn chan DataPoint
	accepts      StreamType
}

// Gather implements Gatherer interface. It records the datapoint
// to the backend, or returns error if the datapoint type is not
// supported by the backend.
func (ddb *dataDogBackend) Gather(dp DataPoint) error {

	switch dp.Type() {
	case GaugeType:
		return ddb.gauge(dp)
	case CountType:
	case EventType:
	case HistogramType:
	case SetType:
	case SimpleEventType:
	default:
		return fmt.Errorf("unsupported datapoint type: %v", dp.Type())
	}

	return nil
}

func (ddb *dataDogBackend) gauge(dp DataPoint) error {

	tags := []string{}
	for k, v := range dp.Tags() {
		tags = append(tags, fmt.Sprintf("%s:%#v", k, v))
	}

	value, ok := dp.Value().(float64)
	if !ok {
		return fmt.Errorf("invalid value type for gauge datapoint")
	}

	return ddb.client.Gauge(dp.MetricName(), value, tags, 1)
}

// Accepts returns whether the backend can handle the
// provided datapoint type.
func (ddb *dataDogBackend) Accepts(typ StreamType) bool {
	return typ&ddb.accepts > 0
}

func NewInfluxBackend() Gatherer {
	// TODO: write stuff here
	return nil
}
