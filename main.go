package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

func runner(ctx context.Context, m *metrics) error {
	var err error
	tags := []string{"foo:bar"}
	var n int64 = 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println("for() loop exiting.")
			return nil
		case <-time.After(time.Duration(n) * time.Second):
			fmt.Printf("golly counter incremented by %d\n", n)
			err = m.statsd.Count("golly", n, tags, float64(n))
			if err != nil {
				log.Fatal(err)
			}
			n = 1 + rand.Int64N(5)
			fmt.Printf("waiting for %d seconds\n", n)
		}
	}
}

func main() {

	fmt.Println("main() starting")

	ctx, cancel := context.WithCancel(context.Background())

	tasks, _ := startTaskManager(cancel)
	metricClient, _ := setupMetrics(tasks)

	tasks.Go(func() error { return runner(ctx, metricClient) })

	<-ctx.Done()
	metricClient.shutdownMetrics()

	fmt.Println("Waiting for all tasks to complete.")
	if err := tasks.Wait(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("main() exiting.")
}
