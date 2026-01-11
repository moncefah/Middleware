package ical

import (
	"bufio"
	"bytes"
	"strings"
)

type Event struct {
	Attributes map[string]string `json:"attributes"`
}

func ParseEvents(rawData []byte) ([]Event, error) {
	scanner := bufio.NewScanner(bytes.NewReader(rawData))

	var events []Event
	currentEvent := Event{Attributes: map[string]string{}}
	currentKey := ""
	currentValue := strings.Builder{}
	inEvent := false

	flushAttribute := func() {
		if currentKey == "" {
			return
		}
		currentEvent.Attributes[currentKey] = currentValue.String()
		currentKey = ""
		currentValue.Reset()
	}

	for scanner.Scan() {
		line := scanner.Text()
		if !inEvent {
			if line == "BEGIN:VEVENT" {
				inEvent = true
				currentEvent = Event{Attributes: map[string]string{}}
			}
			continue
		}

		if line == "END:VEVENT" {
			flushAttribute()
			events = append(events, currentEvent)
			inEvent = false
			continue
		}

		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			currentValue.WriteString(strings.TrimLeft(line, " \t"))
			continue
		}

		flushAttribute()

		splitted := strings.SplitN(line, ":", 2)
		if len(splitted) != 2 {
			continue
		}

		currentKey = splitted[0]
		currentValue.WriteString(splitted[1])
	}

	flushAttribute()
	if inEvent {
		events = append(events, currentEvent)
	}

	return events, scanner.Err()
}
