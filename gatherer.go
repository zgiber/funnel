package funnel

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type Gatherer interface {
	Gather(DataPoint) error
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
	}

	return ddb, nil
}

type dataDogBackend struct {
	client       *statsd.Client
	dataPointsIn chan DataPoint
}

// Gather implements Gatherer interface. It records the datapoint
// to the backend, or returns error if the datapoint type is not
// supported by the backend.
func (ddb *dataDogBackend) Gather(dp DataPoint) error {
	return ddb.gauge(dp)
}

func (ddb *dataDogBackend) gauge(dp DataPoint) error {

	tags := []string{}
	tagsMap := dp.Tags()
	if u := dp.Unit(); u != "" {
		tagsMap["unit"] = u
	}

	for k, v := range dp.Tags() {
		tags = append(tags, fmt.Sprintf("%s:%#v", k, v))
	}

	var value float64
	switch t := dp.Value().(type) {
	case float64:
		value = t
	case float32:
		value = float64(t)
	case int32:
		value = float64(t)
	case int64:
		value = float64(t)
	default:
		return fmt.Errorf("unsupported datapoint value type")
	}

	return ddb.client.Gauge(dp.MetricName(), value, tags, 1)
}

func NewInfluxBackend() Gatherer {
	// TODO: write stuff here
	return nil
}
