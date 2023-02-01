package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
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

	log.Warn().Msg("test warning")
	log.Debug().Msg("test debug")
	log.Error().Msg("test error")
	log.Info().Msgf("Serving request %s %s\n", r.Method, r.URL.Path)
	fmt.Fprintln(w, getenv("WWW_BODY", "Hello, World!"))
}

func StartAPI() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal().Msg("Required value PORT not found in environment, exiting")
	}

	http.HandleFunc("/", handleRequest)
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msgf("Starting server on port %s", port)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msgf("Server is ready to handle requests at %s", port)
}
