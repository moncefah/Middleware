package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"middleware/scheduler/internal/ical"
)

const (
	streamName  = "EVENTS"
	subjectName = "EVENTS.create"
)

type Publisher struct {
	nc  *nats.Conn
	jsc nats.JetStreamContext
}

func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	jsc, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("init jetstream: %w", err)
	}

	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{streamName + ".>"},
	})
	if err != nil && !errors.Is(err, nats.ErrStreamNameAlreadyInUse) && !errors.Is(err, nats.ErrStreamAlreadyExists) {
		nc.Close()
		return nil, fmt.Errorf("add stream: %w", err)
	}

	return &Publisher{
		nc:  nc,
		jsc: jsc,
	}, nil
}

func (p *Publisher) PublishEvent(ctx context.Context, event ical.Event) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	pubAckFuture, err := p.jsc.PublishAsync(subjectName, messageBytes)
	if err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	select {
	case <-pubAckFuture.Ok():
		return nil
	case err := <-pubAckFuture.Err():
		return fmt.Errorf("publish ack error: %w", err)
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(5 * time.Second):
		return errors.New("publish ack timeout")
	}
}

func (p *Publisher) Close() {
	if p.nc != nil {
		p.nc.Close()
	}
}
