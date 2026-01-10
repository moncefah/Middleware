package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"middleware/scheduler/internal/config"
	"middleware/scheduler/internal/ical"
	"middleware/scheduler/internal/mq"
)

const defaultICSURL = "https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp"

func RunOnce(ctx context.Context, client *config.Client, publisher *mq.Publisher) error {
	agendas, err := client.FetchAgendas(ctx)
	if err != nil {
		return err
	}

	resourceIDs := make([]string, 0, len(agendas))
	for _, agenda := range agendas {
		if agenda.UcaID == "" {
			continue
		}
		resourceIDs = append(resourceIDs, agenda.UcaID)
	}

	if len(resourceIDs) == 0 {
		return errors.New("no agenda ids available")
	}

	icalURL := fmt.Sprintf("%s?resources=%s&projectId=3&calType=ical&nbWeeks=4&displayConfigId=128", defaultICSURL, strings.Join(resourceIDs, ","))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, icalURL, nil)
	if err != nil {
		return fmt.Errorf("create ical request: %w", err)
	}

	httpClient := &http.Client{Timeout: 20 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch ical: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected ical status code: %d", resp.StatusCode)
	}

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read ical response: %w", err)
	}

	events, err := ical.ParseEvents(rawData)
	if err != nil {
		return fmt.Errorf("parse ical: %w", err)
	}

	for _, event := range events {
		if err := publisher.PublishEvent(ctx, event); err != nil {
			return err
		}
	}

	payload, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal events: %w", err)
	}

	fmt.Println(string(payload))
	return nil
}
