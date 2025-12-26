package events

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"middleware/go-timetable/internal/helpers"
	"middleware/go-timetable/internal/models"
	"net/http"
)

func Context(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eventId := chi.URLParam(r, "id")
		if eventId == "" {
			body, status := helpers.RespondError(&models.ErrorUnprocessableEntity{
				Message: fmt.Sprintf("cannot parse id (%s)", chi.URLParam(r, "id"))})

			w.WriteHeader(status)
			if body != nil {
				_, _ = w.Write(body)
			}
			return
		}

		ctx := context.WithValue(r.Context(), "eventId", eventId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
