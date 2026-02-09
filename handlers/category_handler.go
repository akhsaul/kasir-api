package handler

import (
	"errors"
	"net/http"

	helper "kasir-api/helpers"
	model "kasir-api/models"
	service "kasir-api/services"
)

// CategoryHandler handles HTTP requests for category endpoints.
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler creates a new instance of CategoryHandler.
func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: svc,
	}
}

// HandleGetAll handles GET /api/categories.
// Supports query parameters: ?page=1&limit=20 for pagination.
func (h *CategoryHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}

	page, limit := helper.ParsePagination(r, 20)
	total := len(categories)

	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	paged := &model.PaginatedResponse{
		Items:      categories[start:end],
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	helper.WriteSuccess(w, http.StatusOK, "Success", paged)
}

// HandleGetByID handles GET /api/categories/{id}.
func (h *CategoryHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/categories/", model.ErrCategoryNotFound)
	if !ok {
		return
	}

	category, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, model.ErrCategoryNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", category)
}

// HandleCreate handles POST /api/categories.
func (h *CategoryHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var category model.Category
	if !helper.ValidatePayload(w, r, &category) {
		return
	}

	createdCategory, err := h.service.Create(&category)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Category created successfully", createdCategory)
}

// HandleUpdate handles PUT /api/categories/{id}.
func (h *CategoryHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/categories/", model.ErrCategoryNotFound)
	if !ok {
		return
	}

	var category model.Category
	if !helper.ValidatePayload(w, r, &category) {
		return
	}

	updatedCategory, err := h.service.Update(id, &category)
	if err != nil {
		if errors.Is(err, model.ErrCategoryNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Category updated successfully", updatedCategory)
}

// HandleDelete handles DELETE /api/categories/{id}.
func (h *CategoryHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/categories/", model.ErrCategoryNotFound)
	if !ok {
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		if errors.Is(err, model.ErrCategoryNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to delete category", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Category deleted successfully", nil)
}
