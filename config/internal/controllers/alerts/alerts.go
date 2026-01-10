package alerts

import (
	"middleware/config/internal/helpers"
	"middleware/config/internal/services/alerts"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetAlerts
// @Tags         alerts
// @Summary      Get alerts for a resource.
// @Description  Get all alerts configured for a specific resource ID.
// @Param        id   path      int  true  "Resource ID"
// @Success      200            {array}  models.Alert
// @Failure      500             "Something went wrong"
// @Router       /resources/{id}/alerts [get]
func GetAlerts(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	res, err := alerts.GetAlertsForResource(id)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	helpers.RespondJSON(w, http.StatusOK, res)
}

// AddAlert
// @Tags         alerts
// @Summary      Add an alert to a resource.
// @Description  Link a mail to a resource to receive alerts.
// @Param        id   path      int  true  "Resource ID"
// @Param        mail_id   path      int  true  "Mail ID"
// @Success      201             "Created"
// @Failure      500             "Something went wrong"
// @Router       /resources/{id}/alerts/{mail_id} [post]
func AddAlert(w http.ResponseWriter, r *http.Request) {
	resIdStr := chi.URLParam(r, "id")
	resId, err := strconv.ParseInt(resIdStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	mailIdStr := chi.URLParam(r, "mail_id")
	mailId, err := strconv.ParseInt(mailIdStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	if err := alerts.AddAlert(mailId, resId); err != nil {
		helpers.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// RemoveAlert
// @Tags         alerts
// @Summary      Remove an alert from a resource.
// @Description  Unlink a mail from a resource.
// @Param        id   path      int  true  "Resource ID"
// @Param        mail_id   path      int  true  "Mail ID"
// @Success      204             "No Content"
// @Failure      500             "Something went wrong"
// @Router       /resources/{id}/alerts/{mail_id} [delete]
func RemoveAlert(w http.ResponseWriter, r *http.Request) {
	resIdStr := chi.URLParam(r, "id")
	resId, err := strconv.ParseInt(resIdStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	mailIdStr := chi.URLParam(r, "mail_id")
	mailId, err := strconv.ParseInt(mailIdStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}

	if err := alerts.RemoveAlert(mailId, resId); err != nil {
		helpers.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
