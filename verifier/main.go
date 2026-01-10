package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Event struct {
	UID         string    `json:"uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// 1. Simulate Publish Scheduler Event (Initial)
	// We need an event that matches the resource name "Test Resource" created in the bash script
	// or matches "arthur.lecomte" if the matching logic compares user?? No, matching logic compares event summary/desc with resource name.

	// Wait a bit for services to be ready
	time.Sleep(2 * time.Second)

	event := Event{
		UID:         "TEST-UID-123456",
		Summary:     "Class for Test Resource", // Contains "Test Resource"
		Description: "A test event",
		Location:    "Room 1",
		Start:       time.Now().Add(1 * time.Hour),
		End:         time.Now().Add(2 * time.Hour),
	}

	// We publish to SCHEDULER.events to populate the DB logic in Config service
	// The config service updates its DB.
	// Ideally we should do this twice. First time "New event detected". Second time "Changed".

	log.Println("Publishing Initial Event...")
	data, _ := json.Marshal([]Event{event})
	nc.Publish("SCHEDULER.events", data)

	time.Sleep(2 * time.Second)

	// 2. Modify Event
	event.Location = "Room 2 (Changed)"
	log.Println("Publishing Modified Event...")
	dataMod, _ := json.Marshal([]Event{event})
	nc.Publish("SCHEDULER.events", dataMod)

	log.Println("Events published. Check Alerter logs for 'Processing modification' and email sending.")
}
