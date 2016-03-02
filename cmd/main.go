package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/zgiber/funnel"
)

var (
	numGoroutinesMetric = funnel.NewMetric("goroutines", "")
	memAllocMetric      = funnel.NewMetric("mem_alloc", "Bytes")
	memTotalMetric      = funnel.NewMetric("mem_total", "Bytes")
	memSysMetric        = funnel.NewMetric("mem_sys", "Bytes")
)

func main() {
	g, err := funnel.NewDatadogGatherer("127.0.0.1:8125", "testapp")
	if err != nil {
		log.Fatal(err)
	}

	m := funnel.NewMetric("test2", "percent")

	randomTicker := time.NewTicker(1 * time.Second).C
	numGoroutinesTicker := time.NewTicker(15 * time.Second).C
	memStatsTicker := time.NewTicker(20 * time.Second).C

	for {
		select {
		case <-randomTicker:
			go g.Gather(random(m))
		case <-numGoroutinesTicker:
			go g.Gather(numGoroutines())
		case <-memStatsTicker:
			go g.GatherBatch(memStats())
		}
	}
}

func random(m *funnel.Metric) funnel.DataPoint {
	return funnel.NewDataPoint(m, rand.Float64(), time.Now().UTC(), map[string]interface{}{"test1": "somevalue", "test2": 23})
}

func numGoroutines() funnel.DataPoint {
	tags := map[string]interface{}{
		"app": "testapp",
	}
	return funnel.NewDataPoint(numGoroutinesMetric, runtime.NumGoroutine(), time.Now().UTC(), tags)
}

func memStats() []funnel.DataPoint {

	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	return []funnel.DataPoint{
		funnel.NewDataPoint(memAllocMetric, ms.Alloc, time.Now().UTC(), nil),
		funnel.NewDataPoint(memTotalMetric, ms.TotalAlloc, time.Now().UTC(), nil),
		funnel.NewDataPoint(memSysMetric, ms.Sys, time.Now().UTC(), nil),
	}
}

// background collector for runtime stuff - in the package
