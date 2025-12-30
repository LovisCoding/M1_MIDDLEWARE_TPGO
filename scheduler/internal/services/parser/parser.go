package parser

import (
	"bufio"
	"io"
	"middleware/scheduler/internal/models"
	"strings"
	"time"
)

const (
	icalTimestampFormat = "20060102T150405Z"
)

func ParseICal(r io.Reader) ([]models.Event, error) {
	scanner := bufio.NewScanner(r)
	var events []models.Event
	var currentEvent map[string]string
	inEvent := false

	var lastKey string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "BEGIN:VEVENT" {
			inEvent = true
			currentEvent = make(map[string]string)
			continue
		}

		if line == "END:VEVENT" {
			if inEvent {
				event, err := mapToEvent(currentEvent)
				if err == nil {
					events = append(events, event)
				}
				inEvent = false
			}
			continue
		}

		if !inEvent {
			continue
		}

		// Handle multi-line data (folded lines start with a space or tab)
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if lastKey != "" {
				currentEvent[lastKey] += strings.TrimPrefix(line, line[:1])
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			lastKey = parts[0]
			currentEvent[lastKey] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func mapToEvent(m map[string]string) (models.Event, error) {
	start, _ := time.Parse(icalTimestampFormat, m["DTSTART"])
	end, _ := time.Parse(icalTimestampFormat, m["DTEND"])

	return models.Event{
		UID:         m["UID"],
		Summary:     m["SUMMARY"],
		Description: m["DESCRIPTION"],
		Location:    m["LOCATION"],
		Start:       start,
		End:         end,
	}, nil
}
