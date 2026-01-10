package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"middleware/scheduler/internal/config"
	"middleware/scheduler/internal/mq"
	"middleware/scheduler/internal/runner"
)

const (
	defaultConfigURL = "http://localhost:8080"
	defaultInterval  = 10 * time.Minute
	intervalEnvVar   = "SCHEDULER_INTERVAL"
	configURLEnvVar  = "CONFIG_BASE_URL"
	natsURLEnvVar    = "NATS_URL"
	runTimeout       = 30 * time.Second
)

func main() {
	configURL := getenv(configURLEnvVar, defaultConfigURL)
	interval := parseInterval(getenv(intervalEnvVar, ""))
	if interval == 0 {
		interval = defaultInterval
	}

	natsURL := getenv(natsURLEnvVar, nats.DefaultURL)

	publisher, err := mq.NewPublisher(natsURL)
	if err != nil {
		log.Fatalf("failed to init nats publisher: %v", err)
	}
	defer publisher.Close()

	client := config.NewClient(configURL)
	ctx := context.Background()

	runOnce := func() {
		runCtx, cancel := context.WithTimeout(ctx, runTimeout)
		defer cancel()
		if err := runner.RunOnce(runCtx, client, publisher); err != nil {
			log.Printf("run failed: %v", err)
		}
	}

	runOnce()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			runOnce()
		case <-quit:
			return
		}
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseInterval(raw string) time.Duration {
	if raw == "" {
		return 0
	}
	interval, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("invalid interval %q: %v", raw, err)
		return 0
	}
	return interval
}
