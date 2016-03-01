package funnel

import "time"

const (
	// Seconds is the format for unix timestamps in a DataPoint
	Seconds = 0
	// Milliseconds is the format for milliseconds timestamps used in certain databases
	Milliseconds = 1
	// Nanoseconds is a format for nanoseconds based timestamps
	Nanoseconds = 2
)

type StreamType int

const (
	GaugeType StreamType = 1 << iota
	CountType
	EventType
	HistogramType
	SetType
	SimpleEventType
)

// DataPoint is an interface representing one measurement or sample
type DataPoint interface {
	MetricName() string
	Type() StreamType
	Tags() map[string]interface{}
	Value() interface{}
	TimeStamp() int64
}

// Metric is a struct with fields that
// have the same value in every DataPoint
// within a measurement/event stream
type Metric struct {
	Name            string
	TimeStampFormat uint8
	Type            StreamType
	Unit            string
}

// NewDataPoint returns a dataPoint for a metric with a given value, time, and tags.
func NewDataPoint(m *Metric, value interface{}, t time.Time, tags map[string]interface{}) DataPoint {
	return &dataPoint{
		m,
		value,
		t,
		tags,
	}
}

type dataPoint struct {
	*Metric
	value interface{}
	t     time.Time
	tags  map[string]interface{}
}

func (dp *dataPoint) MetricName() string {
	return dp.Name
}

func (dp *dataPoint) Tags() map[string]interface{} {
	return dp.tags
}

func (dp *dataPoint) Value() interface{} {
	return dp.value
}

func (dp *dataPoint) Type() StreamType {
	return dp.Metric.Type
}

func (dp *dataPoint) TimeStamp() int64 {
	switch dp.TimeStampFormat {
	case 0:
		return dp.t.UTC().Unix()
	case 1:
		return dp.t.UTC().UnixNano() / 1000
	case 2:
		return dp.t.UTC().UnixNano()
	default:
		// for now we're using seconds as default timestamp
		return dp.t.UTC().Unix()
	}
}
