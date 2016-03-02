package funnel

import "time"

const (
	// Seconds is the format for unix timestamps in a DataPoint
	Seconds = "s"
	// Milliseconds is the format for milliseconds timestamps used in certain databases
	Milliseconds = "ms"
	// Nanoseconds is a format for nanoseconds based timestamps
	Nanoseconds = "s"
)

// DataPoint is an interface representing one measurement or sample
type DataPoint interface {
	MetricName() string
	Tags() map[string]interface{}
	Value() interface{}
	Unit() string
	TimeStamp() int64
}

// NewDataPoint returns a dataPoint for a metric with a given value, time, and tags.
func NewDataPoint(m *Metric, value interface{}, t time.Time, tags map[string]interface{}) DataPoint {
	return &dataPoint{
		Metric: m,
		value:  value,
		t:      t,
		tags:   tags,
	}
}

// Metric is a struct with fields that
// have the same value in every DataPoint
// within a measurement/event stream
type Metric struct {
	name            string
	timeStampFormat string
	unit            string
}

// NewMetric returns a metric which can be used for creating datapoints.
func NewMetric(name, unit, timeStampFmt string) *Metric {
	return &Metric{
		name:            name,
		timeStampFormat: timeStampFmt,
		unit:            unit,
	}
}

type dataPoint struct {
	*Metric
	value interface{}
	t     time.Time
	tags  map[string]interface{}
}

func (dp *dataPoint) MetricName() string {
	return dp.name
}

func (dp *dataPoint) Tags() map[string]interface{} {
	return dp.tags
}

func (dp *dataPoint) Value() interface{} {
	return dp.value
}

func (dp *dataPoint) Unit() string {
	return dp.unit
}

func (dp *dataPoint) TimeStamp() int64 {
	switch dp.timeStampFormat {
	case "s":
		return dp.t.UTC().Unix()
	case "ms":
		return dp.t.UTC().UnixNano() / 1000
	case "ns":
		return dp.t.UTC().UnixNano()
	default:
		// for now we're using seconds as default timestamp
		return dp.t.UTC().Unix()
	}
}
