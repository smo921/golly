package main

import (
	"fmt"
	"log"
	"time"
)

func runner(m *metrics, shutdown <-chan struct{}) error {
	var err error
	tags := []string{"foo:bar"}
	for {
		select {
		case <-shutdown:
			fmt.Println("for() loop exiting.")
			return nil
		default:
			err = m.statsd.Count("golly", 1, tags, 1)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("golly counter incremented")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {

	fmt.Println("main() starting")

	shutdown := make(chan struct{})
	tasks, _ := startTaskManager(shutdown)
	metricClient, _ := setupMetrics(tasks)

	tasks.Go(func() error { return runner(metricClient, shutdown) })

	<-shutdown
	metricClient.shutdownMetrics()
	fmt.Println("Waiting for all tasks to complete.")
	if err := tasks.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("main() exiting.")
}
