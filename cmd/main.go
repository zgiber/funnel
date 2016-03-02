package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/zgiber/funnel"
)

func main() {
	g, err := funnel.NewDataDogBackend("127.0.0.1:8125", "testapp")
	if err != nil {
		log.Fatal(err)
	}

	m := funnel.NewMetric("test2", "percent", "ms")

	for {
		time.Sleep(500 * time.Millisecond)
		g.Gather(random(m))
	}
}

func random(m *funnel.Metric) funnel.DataPoint {
	return funnel.NewDataPoint(m, rand.Float64(), time.Now().UTC(), map[string]interface{}{"test1": "somevalue", "test2": 23})
}
