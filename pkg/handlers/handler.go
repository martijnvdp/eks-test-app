package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eks_test_app_request_count",
			Help: "Total number of connections processed by the server",
		},
		[]string{"method", "route", "status_code"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "eks_test_app_request_duration_seconds",
			Help:    "Time taken for each request to complete",
			Buckets: []float64{0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "route", "status_code"},
	)
)

func getenv(key, defaultString string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultString
	}
	return value
}

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

	datetime := time.Now().Format("2006-01-02 15:04:05.000000000") + ": "
	log.SetPrefix("WARN: " + datetime)
	log.Println("test warning")
	log.SetPrefix("DEBUG: " + datetime)
	log.Println("test debug")
	log.SetPrefix("ERROR: " + datetime)
	log.Println("test error")
	log.SetPrefix("INFO: " + datetime)
	log.Printf("Serving request %s %s\n", r.Method, r.URL.Path)
	fmt.Fprintln(w, getenv("WWW_BODY", "Hello, World!"))
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
