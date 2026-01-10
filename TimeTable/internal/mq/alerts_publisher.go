package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/moncefah/TimeTableAlerter/internal/models"
	"github.com/nats-io/nats.go"
)

const (
	alertStreamName  = "ALERTS"
	alertSubjectName = "ALERTS.create"
)

type AlertsPublisher struct {
	nc  *nats.Conn
	jsc nats.JetStreamContext
}

func NewAlertsPublisher(nc *nats.Conn) (*AlertsPublisher, error) {
	jsc, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("init jetstream: %w", err)
	}

	_, err = jsc.AddStream(&nats.StreamConfig{
		Name:     alertStreamName,
		Subjects: []string{alertStreamName + ".>"},
	})
	if err != nil && !errors.Is(err, nats.ErrStreamNameAlreadyInUse) {
		return nil, fmt.Errorf("add stream: %w", err)
	}

	return &AlertsPublisher{
		nc:  nc,
		jsc: jsc,
	}, nil
}

func (p *AlertsPublisher) PublishChange(ctx context.Context, change models.EventChange) error {
	messageBytes, err := json.Marshal(change)
	if err != nil {
		return fmt.Errorf("marshal change: %w", err)
	}

	pubAckFuture, err := p.jsc.PublishAsync(alertSubjectName, messageBytes)
	if err != nil {
		return fmt.Errorf("publish change: %w", err)
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
