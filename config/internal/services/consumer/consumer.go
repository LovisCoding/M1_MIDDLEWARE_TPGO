package consumer

import (
	"context"
	"encoding/json"
	"time"

	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	eventsRepo "middleware/config/internal/repositories/events"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

func EventConsumer() (jetstream.Consumer, error) {
	// create jetstream context from nats connection
	js, err := jetstream.New(helpers.NatsConn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// get existing stream handle
	stream, err := js.Stream(ctx, "SCHEDULER")
	if err != nil {
		return nil, err
	}

	// getting durable consumer
	consumer, err := stream.Consumer(ctx, "EVENTS_CONSUMER")
	if err != nil {
		// if doesn't exist, create durable consumer
		consumer, err = stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
			Durable:       "EVENTS_CONSUMER",
			Name:          "EVENTS_CONSUMER",
			Description:   "Consumer for scheduler events",
			FilterSubject: "SCHEDULER.events",
		})
		if err != nil {
			return nil, err
		}
		logrus.Infof("Created consumer EVENTS_CONSUMER")
	} else {
		logrus.Infof("Got existing consumer EVENTS_CONSUMER")
	}

	return consumer, nil
}

func Consume(consumer jetstream.Consumer) error {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		var events []models.Event
		if err := json.Unmarshal(msg.Data(), &events); err != nil {
			logrus.Errorf("Error unmarshaling events: %v", err)
			_ = msg.Ack()
			return
		}

		for _, event := range events {
			logrus.Infof("Processing UID: %s (%s)", event.UID, event.Summary)
			if event.UID == "ADE60323032352d323032362d5543412d35353335332d302d32" {
				logrus.Infof("DEBUG: Processing target event. New Location: %s", event.Location)
				existing, _ := eventsRepo.GetByUID(event.UID)
				if existing != nil {
					logrus.Infof("DEBUG: DB Location: %s", existing.Location)
					logrus.Infof("DEBUG: hasChanged: %v", hasChanged(*existing, event))
				} else {
					logrus.Infof("DEBUG: Event not in DB yet")
				}
			}
			// Check for changes
			existingEvent, err := eventsRepo.GetByUID(event.UID)
			if err != nil {
				logrus.Errorf("Error getting event %s: %v", event.UID, err)
				continue
			}

			if existingEvent != nil {
				if hasChanged(*existingEvent, event) {
					logrus.Infof("Event %s has changed. Triggering alert.", event.UID)

					// Publish modification event
					modData, err := json.Marshal(event)
					if err != nil {
						logrus.Errorf("Error marshaling modification event: %v", err)
						continue
					}

					if err := helpers.NatsConn.Publish("SCHEDULER.modification", modData); err != nil {
						logrus.Errorf("Error publishing modification event: %v", err)
					}
				}
			} else {
				logrus.Infof("New event detected: %s", event.UID)
			}

			// Update/Insert in DB
			if err := eventsRepo.Upsert(event); err != nil {
				logrus.Errorf("Error upserting event %s: %v", event.UID, err)
			}
		}

		_ = msg.Ack()
	})

	if err != nil {
		return err
	}

	// Keep running until consumer connection is closed
	<-cc.Closed()
	cc.Stop()

	return nil
}

func hasChanged(old, new models.Event) bool {
	return old.Summary != new.Summary ||
		old.Description != new.Description ||
		old.Location != new.Location ||
		!old.Start.Equal(new.Start) ||
		!old.End.Equal(new.End)
}
