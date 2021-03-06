package background

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zgiber/funnel"
)

type mockGatherer struct {
	sync.RWMutex
	Datapoints map[string]int
}

func (mg *mockGatherer) Gather(dp funnel.Datapoint) error {
	mg.Lock()
	mg.Datapoints[dp.MetricName()]++
	// fmt.Println(dp.MetricName(), dp.Value()) // some visual feedback
	mg.Unlock()
	return nil
}

func (mg *mockGatherer) GatherBatch(dps []funnel.Datapoint) error {
	for _, dp := range dps {
		mg.Gather(dp)
	}

	return nil
}

func TestDatapointsGathered(t *testing.T) {
	mg := &mockGatherer{
		Datapoints: map[string]int{},
	}

	NewCollector(mg, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	assert.NotZero(t, mg.Datapoints["goroutines"])
	assert.NotZero(t, mg.Datapoints["mem_alloc"])
	assert.NotZero(t, mg.Datapoints["mem_total"])
	assert.NotZero(t, mg.Datapoints["mem_sys"])
}
