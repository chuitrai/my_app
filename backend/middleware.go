package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_http_request_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path", "method", "code"}) // Phân loại theo path, method, status code

	totalRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "code"},
	)
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Ghi lại response code, vì response writer gốc không làm điều đó
		rec := &statusRecorder{ResponseWriter: w, Status: http.StatusOK}

		// Gọi handler tiếp theo trong chuỗi
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		// Lấy route template để tránh tạo quá nhiều label khác nhau
		// Ví dụ: /api/users/1 và /api/users/2 sẽ đều được gom vào /api/users/{id}
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		statusCode := strconv.Itoa(rec.Status)

		// Ghi nhận metrics
		httpDuration.WithLabelValues(path, r.Method, statusCode).Observe(duration.Seconds())
		totalRequests.WithLabelValues(path, r.Method, statusCode).Inc()
	})
}

// Helper để ghi lại status code
type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
