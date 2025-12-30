package events

import (
	"encoding/json"
	"middleware/go-timetable/internal/helpers"
	"middleware/go-timetable/internal/services/events"
	"net/http"
)

// GetEvents
// @Tags         events
// @Summary      Get all events.
// @Description  Get all events from the timetable, optionally filtered by resourceId.
// @Param        resourceId     query     string  false  "Resource ID to filter events"
// @Success      200            {array}  models.Event
// @Failure      500            "Something went wrong"
// @Router       /events [get]
func GetEvents(w http.ResponseWriter, r *http.Request) {
	// Récupération du paramètre resourceId dans l'URL (?resourceId=...)
	resourceId := r.URL.Query().Get("resourceId")

	eventsList, err := events.GetAllEvents(resourceId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(eventsList)
	_, _ = w.Write(body)
}