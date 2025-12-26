package events

import (
	"encoding/json"
	"middleware/go-timetable/internal/helpers"
	"middleware/go-timetable/internal/services/events"
	"net/http"
)

// GetEvent
// @Tags         events
// @Summary      Get an event.
// @Description  Get a specific event by ID.
// @Param        id             path      string  true  "Event ID"
// @Success      200            {object}  models.Event
// @Failure      422            "Cannot parse id"
// @Failure      500            "Something went wrong"
// @Router       /events/{id} [get]
func GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	eventId, _ := ctx.Value("eventId").(string)

	event, err := events.GetEventById(eventId)
	if err != nil {
		body, status := helpers.RespondError(err)
		w.WriteHeader(status)
		if body != nil {
			_, _ = w.Write(body)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	body, _ := json.Marshal(event)
	_, _ = w.Write(body)
}
