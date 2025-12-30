package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"middleware/scheduler/internal/services/config"
	"middleware/scheduler/internal/services/parser"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	natsService "middleware/scheduler/internal/services/nats"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/zhashkevych/scheduler"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables or defaults")
	}

	baseURL := os.Getenv("CONFIG_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	// Init NATS
	natsClient, err := natsService.NewClient(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer natsClient.Close()

	// Init Stream
	err = natsClient.InitStream("SCHEDULER", []string{"SCHEDULER.events"})
	if err != nil {
		log.Printf("Error initializing stream: %v", err)
	}

	// Define the job
	job := func(ctx context.Context) {
		log.Println("Starting scheduled job...")

		// 1. Fetch resource IDs
		configAPIURL := baseURL + "resources"
		ids, err := config.FetchResourceIDs(configAPIURL)
		if err != nil || len(ids) == 0 {
			if err != nil {
				log.Printf("Warning: could not fetch resource IDs from Config API: %v. Using defaults.", err)
			}
			ids = []int64{13295, 13345}
		}

		// 2. Construct EDT iCal URL
		strIDs := make([]string, len(ids))
		for i, id := range ids {
			strIDs[i] = fmt.Sprintf("%d", id)
		}
		resourcesParam := strings.Join(strIDs, ",")
		edtURL := fmt.Sprintf("https://edt.uca.fr/jsp/custom/modules/plannings/anonymous_cal.jsp?resources=%s&projectId=3&calType=ical&nbWeeks=4&displayConfigId=128", resourcesParam)

		// 3. Fetch iCal data
		resp, err := http.Get(edtURL)
		if err != nil {
			log.Printf("Error fetching iCal: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Unexpected status code from EDT: %d", resp.StatusCode)
			return
		}

		// 4. Parse iCal data
		events, err := parser.ParseICal(resp.Body)
		if err != nil {
			log.Printf("Error parsing iCal: %v", err)
			return
		}

		// 5. Publish events to NATS
		// Chunking to avoid "maximum payload exceeded"
		chunkSize := 50
		for i := 0; i < len(events); i += chunkSize {
			end := i + chunkSize
			if end > len(events) {
				end = len(events)
			}
			chunk := events[i:end]

			data, err := json.Marshal(chunk)
			if err != nil {
				log.Printf("Error marshaling chunk: %v", err)
				continue
			}

			if err := natsClient.Publish("SCHEDULER.events", data); err != nil {
				log.Printf("Error publishing chunk to NATS: %v", err)
				continue
			}
		}

		log.Printf("Successfully published %d events to NATS (in chunks)", len(events))
	}

	// Init Scheduler
	ctx := context.Background()
	s := scheduler.NewScheduler()
	s.Add(ctx, job, time.Second*20) // Run every 20 seconds for testing

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	s.Stop()
}
