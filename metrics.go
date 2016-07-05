package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var metrics = struct {
	iterations    prometheus.Counter
	errors        prometheus.Counter
	awsApiLatency *prometheus.HistogramVec
}{
	prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "databricks",
			Subsystem: "spot_tagger",
			Name:      "iterations",
			Help:      "Number of times this has run",
		},
	),
	prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "databricks",
			Subsystem: "spot_tagger",
			Name:      "errors",
			Help:      "Number of errors that have occurred",
		},
	),
	prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "databricks",
			Subsystem: "spot_tagger",
			Name:      "aws_api_latlency",
			Help:      "Seconds spent waiting for AWS API Calls",
		},
		[]string{"apiCall"},
	),
}

func init() {
	prometheus.MustRegister(metrics.iterations)
	prometheus.MustRegister(metrics.errors)
	prometheus.MustRegister(metrics.awsApiLatency)

	go func() {
		for {
			http.Handle("/metrics", prometheus.Handler())
			if err := http.ListenAndServe(":8080", nil); err != nil {
				log.Println(errors.Wrap(err, "failed to start http server"))
			}
		}
	}()
}
