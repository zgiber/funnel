package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/zgiber/funnel"
	"github.com/zgiber/funnel/background"
)

func main() {
	g, err := funnel.NewDatadogGatherer("127.0.0.1:8125", "testapp")
	if err != nil {
		log.Fatal(err)
	}

	// collecting background stats
	background.NewCollector(g, 15*time.Second)

	// collecting custom stats
	m := funnel.NewMetric("test2", "percent")
	randomTicker := time.NewTicker(1 * time.Second).C

	for {
		select {
		case <-randomTicker:
			go g.Gather(random(m))
		}
	}
}

func random(m *funnel.Metric) funnel.Datapoint {
	return funnel.NewDatapoint(m, rand.Float64(), time.Now().UTC(), map[string]interface{}{"test1": "somevalue", "test2": 23})
}
