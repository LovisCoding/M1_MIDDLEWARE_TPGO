package resources

import (
	"encoding/json"
	"middleware/config/internal/helpers"
	"middleware/config/internal/models"
	"middleware/config/internal/services/resources"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetResources
// @Tags         resources
// @Summary      Get all resources.
// @Description  Get all resources.
// @Success      200            {array}  models.Ressource
// @Failure      500             "Something went wrong"
// @Router       /resources [get]
func GetResources(w http.ResponseWriter, _ *http.Request) {
	res, err := resources.GetAllResources()
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	helpers.RespondJSON(w, http.StatusOK, res)
}

// GetResource
// @Tags         resources
// @Summary      Get resource by ID.
// @Description  Get resource by ID.
// @Param        id   path      int  true  "Resource ID"
// @Success      200            {object}  models.Ressource
// @Failure      500             "Something went wrong"
// @Router       /resources/{id} [get]
func GetResource(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid ID"})
		return
	}

	res, err := resources.GetResourceById(id)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	helpers.RespondJSON(w, http.StatusOK, res)
}

// CreateResource
// @Tags         resources
// @Summary      Create a resource.
// @Description  Create a resource.
// @Param        resource   body      models.Ressource  true  "Resource content"
// @Success      201            {object}  models.Ressource
// @Failure      500             "Something went wrong"
// @Router       /resources [post]
func CreateResource(w http.ResponseWriter, r *http.Request) {
	var req models.Ressource
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid Request Body"})
		return
	}

	id, err := resources.CreateResource(req.Name, req.Type)
	if err != nil {
		helpers.RespondError(w, err)
		return
	}
	req.ID = id
	helpers.RespondJSON(w, http.StatusCreated, req)
}

// UpdateResource
// @Tags         resources
// @Summary      Update a resource.
// @Description  Update a resource.
// @Param        id   path      int  true  "Resource ID"
// @Param        resource   body      models.Ressource  true  "Resource content"
// @Success      200            {object}  models.Ressource
// @Failure      500             "Something went wrong"
// @Router       /resources/{id} [put]
func UpdateResource(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid ID"})
		return
	}

	var req models.Ressource
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid Request Body"})
		return
	}

	if err := resources.UpdateResource(id, req.Name, req.Type); err != nil {
		helpers.RespondError(w, err)
		return
	}
	req.ID = id
	helpers.RespondJSON(w, http.StatusOK, req)
}

// DeleteResource
// @Tags         resources
// @Summary      Delete a resource.
// @Description  Delete a resource.
// @Param        id   path      int  true  "Resource ID"
// @Success      204             "No Content"
// @Failure      500             "Something went wrong"
// @Router       /resources/{id} [delete]
func DeleteResource(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		helpers.RespondError(w, &models.ErrorUnprocessableEntity{Message: "Invalid ID"})
		return
	}

	if err := resources.DeleteResource(id); err != nil {
		helpers.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
