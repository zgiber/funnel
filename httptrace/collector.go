package httptrace

import (
	"net/http"
	"time"

	"github.com/zgiber/funnel"
)

var (
	httpMetric = funnel.NewMetric("http_trace", "ms")
)

type httpWriter struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
}

func (w *httpWriter) WriteHeader(status int) {
	if !w.wroteHeader {
		w.status = status
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w httpWriter) Status() int {
	return w.status
}

func New(g funnel.Gatherer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			hw := &httpWriter{ResponseWriter: w}

			t1 := time.Now()
			h.ServeHTTP(hw, r)

			if hw.Status() == 0 {
				hw.WriteHeader(http.StatusOK)
			}
			// t2 := time.Now()

			go func() {
				// create new datapoint
				tags := map[string]interface{}{
					"status": hw.Status(),
					"path":   r.URL.Path,
				}

				datapoint := funnel.NewDatapoint(httpMetric, time.Since(t1), time.Now().UTC(), tags)
				g.Gather(datapoint)
			}()
		}

		return http.HandlerFunc(fn)
	}
}
