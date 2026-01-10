package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/moncefah/TimeTableAlerter/internal/helpers"
	"github.com/moncefah/TimeTableAlerter/internal/models"
	"github.com/moncefah/TimeTableAlerter/internal/mq"
	repository "github.com/moncefah/TimeTableAlerter/internal/repositories/events"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

const (
	streamName    = "EVENTS"
	consumerName  = "timetable_consumer"
	subjectFilter = "EVENTS.create"
)

type schedulerEvent struct {
	Attributes map[string]string `json:"attributes"`
}

func EventConsumer() (*jetstream.Consumer, error) {
	if helpers.NatsConn == nil {
		return nil, errors.New("nats connection not initialized")
	}

	js, err := jetstream.New(helpers.NatsConn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		return nil, err
	}

	consumer, err := stream.Consumer(ctx, consumerName)
	if err != nil {
		consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:       consumerName,
			Name:          consumerName,
			Description:   "Timetable events consumer",
			FilterSubject: subjectFilter,
			AckPolicy:     jetstream.AckExplicitPolicy,
			ReplayPolicy:  jetstream.ReplayInstantPolicy,
			DeliverPolicy: jetstream.DeliverAllPolicy,
			MaxAckPending: 256,
		})
		if err != nil {
			return nil, err
		}
		logrus.Infof("Created consumer")
	} else {
		logrus.Infof("Got existing consumer")
	}

	return &consumer, nil
}

func Consume(consumer jetstream.Consumer) error {
	publisher, err := mq.NewAlertsPublisher(helpers.NatsConn)
	if err != nil {
		return err
	}

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		if err := handleMessage(msg.Data(), publisher); err != nil {
			logrus.Warnf("failed to handle message: %v", err)
			return
		}
		if err := msg.Ack(); err != nil {
			logrus.Warnf("ack failed: %v", err)
		}
	})
	if err != nil {
		return err
	}

	<-cc.Closed()
	cc.Stop()

	return nil
}

func handleMessage(data []byte, publisher *mq.AlertsPublisher) error {
	var payload schedulerEvent
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("decode event: %w", err)
	}

	event, err := toEvent(payload.Attributes)
	if err != nil {
		return err
	}

	existing, err := repository.GetEventByUID(event.UID)
	if err != nil {
		return err
	}

	if existing == nil {
		event.ID = uuid.Must(uuid.NewV4())
		return repository.CreateEvent(event)
	}

	changes := diffEvent(*existing, event)
	if len(changes) == 0 {
		return nil
	}

	event.ID = existing.ID
	if err := repository.UpdateEvent(event); err != nil {
		return err
	}

	change := models.EventChange{
		EventID:    event.ID.String(),
		UID:        event.UID,
		Summary:    event.Name,
		Start:      event.Start.Format(time.RFC3339),
		End:        event.End.Format(time.RFC3339),
		Location:   event.Location,
		ChangeType: "updated",
		Changes:    changes,
	}

	return publisher.PublishChange(context.Background(), change)
}

func toEvent(attributes map[string]string) (models.Event, error) {
	start, err := parseICalTime(attributes["DTSTART"])
	if err != nil {
		return models.Event{}, fmt.Errorf("parse DTSTART: %w", err)
	}
	end, err := parseICalTime(attributes["DTEND"])
	if err != nil {
		return models.Event{}, fmt.Errorf("parse DTEND: %w", err)
	}
	lastUpdate, err := parseICalTime(attributes["LAST-MODIFIED"])
	if err != nil {
		lastUpdate, err = parseICalTime(attributes["DTSTAMP"])
		if err != nil {
			lastUpdate = time.Now().UTC()
		}
	}

	return models.Event{
		AgendaIDs:   []string{},
		UID:         attributes["UID"],
		Description: attributes["DESCRIPTION"],
		Name:        attributes["SUMMARY"],
		Start:       start,
		End:         end,
		Location:    attributes["LOCATION"],
		LastUpdate:  lastUpdate,
	}, nil
}

func parseICalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("empty time")
	}
	if parsed, err := time.Parse("20060102T150405Z", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("20060102T150405", value); err == nil {
		return parsed, nil
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", value)
}

func diffEvent(existing models.Event, incoming models.Event) map[string][2]string {
	changes := make(map[string][2]string)
	if existing.Name != incoming.Name {
		changes["summary"] = [2]string{existing.Name, incoming.Name}
	}
	if existing.Description != incoming.Description {
		changes["description"] = [2]string{existing.Description, incoming.Description}
	}
	if !existing.Start.Equal(incoming.Start) {
		changes["start"] = [2]string{existing.Start.Format(time.RFC3339), incoming.Start.Format(time.RFC3339)}
	}
	if !existing.End.Equal(incoming.End) {
		changes["end"] = [2]string{existing.End.Format(time.RFC3339), incoming.End.Format(time.RFC3339)}
	}
	if existing.Location != incoming.Location {
		changes["location"] = [2]string{existing.Location, incoming.Location}
	}

	return changes
}
