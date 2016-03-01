package funnel

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
)

type Gatherer interface {
	Gather(dp DataPoint) error
}

func NewDataDogBackend(addr string) (Gatherer, error) {
	var c *statsd.Client
	var err error

	c, err = statsd.NewBuffered(addr, 128)
	if err != nil {
		return nil, err
	}

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

func (ddb *dataDogBackend) gauge(dp DataPoint) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(fmt.Sprint(r))
		}
	}()

	tags := []string{}
	for k, v := range dp.Tags() {
		tags = append(tags, fmt.Sprintf("%s:%#v", k, v))
	}

	value := dp.Value().(float64)

	return ddb.client.Gauge(dp.MetricName(), value, tags, 1)
}

func (ddb *dataDogBackend) canHandle(typ StreamType) bool {
	return typ&ddb.accepts > 0
}

func NewInfluxBackend() Gatherer {
	// TODO: write stuff here
	return nil
}
