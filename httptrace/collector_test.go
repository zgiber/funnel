package httptrace

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

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
	fmt.Println(dp.MetricName(), dp.Value()) // some visual feedback
	mg.Unlock()
	return nil
}

func (mg *mockGatherer) GatherBatch(dps []funnel.Datapoint) error {
	for _, dp := range dps {
		mg.Gather(dp)
	}

	return nil
}

func TestMiddlewareTrackStatusAndUrl(t *testing.T) {
	g := &mockGatherer{
		Datapoints: map[string]int{},
	}

	hf := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	}

	mw := New(g)

	ts := httptest.NewServer(mw(http.HandlerFunc(hf)))
	http.Get(ts.URL + "/somepath")

	assert.NotZero(t, g.Datapoints["http_trace"])

}
