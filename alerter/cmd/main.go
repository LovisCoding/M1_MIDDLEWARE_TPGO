package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"time"

	"middleware/alerter/internal/models"
	"middleware/alerter/internal/services/config"
	"middleware/alerter/internal/services/mailer"
	"middleware/alerter/internal/services/templater"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

func main() {
	// NATS Connection
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		logrus.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		logrus.Fatalf("Error creating JetStream context: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	stream, err := js.Stream(ctx, "SCHEDULER")
	if err != nil {
		logrus.Fatalf("Error getting stream: %v", err)
	}

	// Create durable consumer for Alerter
	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       "ALERTER_SERVICE",
		Name:          "ALERTER_SERVICE",
		Description:   "Consumer for event modifications",
		FilterSubject: "SCHEDULER.modification",
	})
	if err != nil {
		logrus.Fatalf("Error creating consumer: %v", err)
	}

	logrus.Info("Alerter service started, listening for modifications...")

	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		var event models.Event
		if err := json.Unmarshal(msg.Data(), &event); err != nil {
			logrus.Errorf("Error unmarshaling event: %v", err)
			_ = msg.Ack()
			return
		}

		logrus.Infof("Processing modification for event: %s (%s)", event.UID, event.Summary)

		// 1. Find relevant resources
		resources, err := config.GetAllResources()
		if err != nil {
			logrus.Errorf("Error fetching resources: %v", err)
			_ = msg.Ack() // Retry? For now ack to avoid loop
			return
		}

		// Naive matching: check if resource name is contained in event summary or description
		// In a real app, we would have a better link (e.g., tags, group IDs)
		var matchedResources []models.Resource
		for _, res := range resources {
			// Basic case-insensitive check could be better
			if strings.Contains(strings.ToLower(event.Summary), strings.ToLower(res.Name)) ||
				strings.Contains(strings.ToLower(event.Description), strings.ToLower(res.Name)) {
				matchedResources = append(matchedResources, res)
			}
		}

		if len(matchedResources) == 0 {
			logrus.Infof("No matching resources found for event %s", event.UID)
			_ = msg.Ack()
			return
		}

		// 2. Find subscribers for these resources
		uniqueEmails := make(map[string]bool)
		for _, res := range matchedResources {
			alerts, err := config.GetAlertsForResource(res.ID)
			if err != nil {
				logrus.Errorf("Error fetching alerts for resource %s: %v", res.Name, err)
				continue
			}

			for _, alert := range alerts {
				mail, err := config.GetMail(alert.MailID)
				if err != nil {
					logrus.Errorf("Error fetching mail %d: %v", alert.MailID, err)
					continue
				}
				uniqueEmails[mail.Mail] = true
			}
		}

		if len(uniqueEmails) == 0 {
			logrus.Info("No subscribers found for affected resources")
			_ = msg.Ack()
			return
		}

		// 3. Prepare Email
		body, matter, err := templater.GetStringFromEmbeddedTemplate("modified_event.html", event)
		if err != nil {
			logrus.Errorf("Error preparing template: %v", err)
			_ = msg.Ack()
			return
		}

		// 4. Send Email
		var recipients []string
		for email := range uniqueEmails {
			recipients = append(recipients, email)
		}

		// The API might accept multiple recipients, but to be safe/personalized per limit,
		// we send one by one or in batch.
		// Let's assume the API handles the list.
		if err := mailer.SendEmail(recipients, matter.Subject, body); err != nil {
			logrus.Errorf("Error sending email: %v", err)
			// Decide if we Nack to retry later
		}

		_ = msg.Ack()
	})

	if err != nil {
		logrus.Fatalf("Error consuming: %v", err)
	}

	// Keep running
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	cc.Stop()
}
