package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func startTaskManager(cancel context.CancelFunc) (*errgroup.Group, error) {
	tasks, _ := errgroup.WithContext(context.Background())

	// SIGINT handler
	tasks.Go(func() error {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		fmt.Println("SIGINT received, shutting down.")

		cancel()
		return nil
	})

	return tasks, nil
}
