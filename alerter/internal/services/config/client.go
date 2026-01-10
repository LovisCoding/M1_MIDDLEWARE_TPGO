package config

import (
	"encoding/json"
	"fmt"
	"middleware/alerter/internal/models"
	"net/http"
	"os"
)

var baseURL string

func init() {
	baseURL = os.Getenv("CONFIG_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}
}

func GetAllResources() ([]models.Resource, error) {
	resp, err := http.Get(baseURL + "resources")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var resources []models.Resource
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, err
	}

	return resources, nil
}

func GetAlertsForResource(resourceID int64) ([]models.Alert, error) {
	url := fmt.Sprintf("%sresources/%d/alerts", baseURL, resourceID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var alerts []models.Alert
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return nil, err
	}

	return alerts, nil
}

func GetMail(mailID int64) (*models.Mail, error) {
	url := fmt.Sprintf("%smails/%d", baseURL, mailID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var mail models.Mail
	if err := json.NewDecoder(resp.Body).Decode(&mail); err != nil {
		return nil, err
	}

	return &mail, nil
}
