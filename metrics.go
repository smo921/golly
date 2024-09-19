package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

type metrics struct {
	statsd     *statsd.Client
	prometheus *http.Server
}

/*func setupOpenTelemetryMetrics() error {
	return errors.New("not implemented")
}*/

func setupPrometheusMetrics(srv *http.Server) error {
	http.Handle("/metrics", promhttp.Handler())

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return fmt.Errorf("HTTP server ListenAndServe: %v", err)
	}

	fmt.Println("HTTP Server exiting.")
	return nil
}

func setupStatsd() (*statsd.Client, error) {
	url := os.Getenv("GOLLY_STATSD_URL")
	if url == "" {
		url = "localhost:8125"
	}

	fmt.Println("Using statsd url:", url)

	statsd, err := statsd.New(url)
	return statsd, err
}

func setupMetrics(tasks *errgroup.Group) (*metrics, error) {
	s, err := setupStatsd()
	if err != nil {
		return nil, err
	}
	p := &http.Server{Addr: ":2112"}
	tasks.Go(func() error { return setupPrometheusMetrics(p) })

	m := &metrics{statsd: s, prometheus: p}
	return m, nil
}

func (m *metrics) shutdownMetrics() error {
	fmt.Println("Http Server shutdown initiated.")
	if err := m.prometheus.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		return fmt.Errorf("HTTP server Shutdown error: %v", err)
	}
	return nil
}
