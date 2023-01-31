package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "Total number of connections processed by the server",
		},
		[]string{"method", "route", "status_code"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Time taken for each request to complete",
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "route", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	duration := time.Since(start).Seconds()
	requestDuration.With(prometheus.Labels{
		"method":      r.Method,
		"route":       r.URL.Path,
		"status_code": "200",
	}).Observe(duration)
	requestCounter.With(prometheus.Labels{
		"method":      r.Method,
		"route":       r.URL.Path,
		"status_code": "200",
	}).Inc()
	log.Printf("Serving request %s %s\n", r.Method, r.URL.Path)
	fmt.Fprintln(w, "Hello, World!")
}

func StartAPI() {
	http.HandleFunc("/", handleRequest)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
