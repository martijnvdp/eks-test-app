package handlers

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	connectionsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "connections_total",
			Help: "Total number of connections processed by the server",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(connectionsTotal)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	connectionsTotal.With(prometheus.Labels{"status": "success"}).Inc()
	fmt.Fprintln(w, "Hello, World!")
}

func StartAPI() {
	http.HandleFunc("/", handleRequest)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
