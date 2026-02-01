package handler

import (
	"errors"
	"kasir-api/models"
	"kasir-api/helpers"
	"kasir-api/services"
	"net/http"
	"strconv"
	"strings"
)

// CategoryHandler handles HTTP requests for category endpoints
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler creates a new instance of CategoryHandler
func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// HandleGetAll handles GET /api/categories
func (h *CategoryHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", categories)
}

// HandleGetByID handles GET /api/categories/{id}
func (h *CategoryHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid Category ID", err)
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, "Category not found", err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", category)
}

// HandleCreate handles POST /api/categories
func (h *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	if !helper.ValidatePayload(w, r, &category) {
		return
	}

	createdCategory, err := h.service.Create(&category)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to create category", err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Category created successfully", createdCategory)
}

// HandleUpdate handles PUT /api/categories/{id}
func (h *CategoryHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid Category ID", err)
		return
	}

	var category model.Category
	if !helper.ValidatePayload(w, r, &category) {
		return
	}

	updatedCategory, err := h.service.Update(id, &category)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, "Category not found", err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to update category", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Category updated successfully", updatedCategory)
}

// HandleDelete handles DELETE /api/categories/{id}
func (h *CategoryHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid Category ID", err)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, "Category not found", err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to delete category", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Category deleted successfully", nil)
}
