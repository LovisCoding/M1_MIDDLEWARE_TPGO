package config

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Resource struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func FetchResourceIDs(apiURL string) ([]int64, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var resources []Resource
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, err
	}

	var ids []int64
	for _, r := range resources {
		ids = append(ids, r.ID)
	}

	return ids, nil
}
