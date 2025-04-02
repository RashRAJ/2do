package metrics_monitoring

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Define metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	tasksCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_created_total",
			Help: "Total number of tasks created.",
		},
	)

	tasksCompleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_completed_total",
			Help: "Total number of tasks completed.",
		},
	)

	tasksDeleted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tasks_deleted_total",
			Help: "Total number of tasks deleted.",
		},
	)

	dbQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Duration of database queries.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query"},
	)

	dbErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors.",
		},
	)
)

// Register metrics with Prometheus
func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(tasksCreated)
	prometheus.MustRegister(tasksCompleted)
	prometheus.MustRegister(tasksDeleted)
	prometheus.MustRegister(dbQueryDuration)
	prometheus.MustRegister(dbErrorsTotal)
}

// InstrumentHTTPRequest instruments an HTTP request with Prometheus metrics
func InstrumentHTTPRequest(method, endpoint string, statusCode int, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, endpoint, http.StatusText(statusCode)).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// IncrementTasksCreated increments the tasks created counter
func IncrementTasksCreated() {
	tasksCreated.Inc()
}

// IncrementTasksCompleted increments the tasks completed counter
func IncrementTasksCompleted() {
	tasksCompleted.Inc()
}

// IncrementTasksDeleted increments the tasks deleted counter
func IncrementTasksDeleted() {
	tasksDeleted.Inc()
}

// InstrumentDBQuery instruments a database query with Prometheus metrics
func InstrumentDBQuery(query string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(query).Observe(duration.Seconds())
}

// IncrementDBErrors increments the database errors counter
func IncrementDBErrors() {
	dbErrorsTotal.Inc()
}

// func ExposeMetrics() {
// 	// Create a new router
// 	router := mux.NewRouter()

// 	// Register the metrics handler
// 	router.Handle("/metrics", promhttp.Handler())

// 	// Start the HTTP server
// 	http.ListenAndServe(":8080", router)
// }
