package mails

import (
	"encoding/json"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"middleware/config/internal/services/mails"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetMails
// @Tags         mails
// @Summary      Get all mails.
// @Description  Get all mails.
// @Success      200            {array}  models.Mail
// @Failure      500             "Something went wrong"
// @Router       /mails [get]
func GetMails(w http.ResponseWriter, _ *http.Request) {
	mailsList, err := mails.GetAllMails()
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	helpers.RespondJSON(w, http.StatusOK, mailsList)
}

// GetMail
// @Tags         mails
// @Summary      Get mail by ID.
// @Description  Get mail by ID.
// @Param        id   path      int  true  "Mail ID"
// @Success      200            {object}  models.Mail
// @Failure      500             "Something went wrong"
// @Router       /mails/{id} [get]
func GetMail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid ID"})
		return
	}

	mail, err := mails.GetMailById(id)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	helpers.RespondJSON(w, http.StatusOK, mail)
}

// CreateMail
// @Tags         mails
// @Summary      Create a mail.
// @Description  Create a mail.
// @Param        mail   body      models.Mail  true  "Mail content"
// @Success      201            {object}  models.Mail
// @Failure      500             "Something went wrong"
// @Router       /mails [post]
func CreateMail(w http.ResponseWriter, r *http.Request) {
	var req models.Mail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid Request Body"})
		return
	}

	id, err := mails.CreateMail(req.Mail)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	req.ID = id
	helpers.RespondJSON(w, http.StatusCreated, req)
}

// UpdateMail
// @Tags         mails
// @Summary      Update a mail.
// @Description  Update a mail.
// @Param        id   path      int  true  "Mail ID"
// @Param        mail   body      models.Mail  true  "Mail content"
// @Success      200            {object}  models.Mail
// @Failure      500             "Something went wrong"
// @Router       /mails/{id} [put]
func UpdateMail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid ID"})
		return
	}

	var req models.Mail
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid Request Body"})
		return
	}

	if err := mails.UpdateMail(id, req.Mail); err != nil {
		helpers.RespondError(w, err)
		return
	}
	req.ID = id
	helpers.RespondJSON(w, http.StatusOK, req)
}

// DeleteMail
// @Tags         mails
// @Summary      Delete a mail.
// @Description  Delete a mail.
// @Param        id   path      int  true  "Mail ID"
// @Success      204             "No Content"
// @Failure      500             "Something went wrong"
// @Router       /mails/{id} [delete]
func DeleteMail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid ID"})
		return
	}

	if err := mails.DeleteMail(id); err != nil {
		helpers.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
