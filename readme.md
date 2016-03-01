# funnel
--
    import "github.com/zgiber/funnel"


## Usage

#### type DataPoint

```go
type DataPoint interface {
	MetricName() string
	Type() StreamType
	Tags() map[string]interface{}
	Value() interface{}
	TimeStamp() int64
}
```

DataPoint is an interface representing one measurement or sample

#### func  NewDataPoint

```go
func NewDataPoint(m *Metric, value interface{}, t time.Time, tags map[string]interface{}) DataPoint
```
NewDataPoint returns a dataPoint for a metric with a given value, time, and
tags.

#### type Gatherer

```go
type Gatherer interface {
	Gather(DataPoint) error
	Accepts(StreamType) bool
}
```


#### func  NewDataDogBackend

```go
func NewDataDogBackend(addr string) (Gatherer, error)
```
NewDataDogBackend returns a Gatherer which collects datapoints to the specified
DataDog address. It uses the datadog statsd package with the buffered client.

#### func  NewInfluxBackend

```go
func NewInfluxBackend() Gatherer
```

#### type Metric

```go
type Metric struct {
	Name            string
	TimeStampFormat uint8
	Type            StreamType
	Unit            string
}
```

Metric is a struct with fields that have the same value in every DataPoint
within a measurement/event stream

#### type StreamType

```go
type StreamType int
```


```go
const (
	// Seconds is the format for unix timestamps in a DataPoint
	Seconds = 0
	// Milliseconds is the format for milliseconds timestamps used in certain databases
	Milliseconds = 1
	// Nanoseconds is a format for nanoseconds based timestamps
	Nanoseconds = 2

	// GaugeType is for recording float64 metric at given times
	GaugeType StreamType = 1 << iota
	// CountType is for recording the number of received entries in one second
	CountType
	// EventType is for recording an event as it happens
	EventType
	// HistogramType tracks the statistical distribution of a set of values
	HistogramType
	// SetType counts the number of unique elements in a group
	SetType
	// SimpleEventType sends an event with the provided title and text
	SimpleEventType
)
```
