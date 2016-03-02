package background

import (
	"log"
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

type Collector struct {
	frequency time.Duration
	gatherer  funnel.BatchGatherer
}

func NewCollector(g funnel.BatchGatherer, freq time.Duration) *Collector {
	c := &Collector{
		frequency: freq,
		gatherer:  g,
	}

	go c.Collect()
	return c
}

func (bgc *Collector) Collect() {
	tick := time.NewTicker(bgc.frequency).C
	for {
		select {
		case <-tick:
			err := bgc.gatherer.GatherBatch(memStats())
			if err != nil {
				log.Println(err)
			}

			err = bgc.gatherer.Gather(numGoroutines())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func numGoroutines() funnel.Datapoint {
	tags := map[string]interface{}{
		"app": "testapp",
	}
	return funnel.NewDatapoint(numGoroutinesMetric, runtime.NumGoroutine(), time.Now().UTC(), tags)
}

func memStats() []funnel.Datapoint {

	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	return []funnel.Datapoint{
		funnel.NewDatapoint(memAllocMetric, ms.Alloc, time.Now().UTC(), nil),
		funnel.NewDatapoint(memTotalMetric, ms.TotalAlloc, time.Now().UTC(), nil),
		funnel.NewDatapoint(memSysMetric, ms.Sys, time.Now().UTC(), nil),
	}
}
