// Package handler: Handler layer — menerima request dan response.
// Error request (invalid ID, invalid JSON, dll) → cek sini.
package handler

import (
	"errors"
	"net/http"

	helper "kasir-api/helpers"
	model "kasir-api/models"
	service "kasir-api/services"
)

// ProductHandler handles HTTP requests for product endpoints.
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler creates a new instance of ProductHandler.
func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: svc,
	}
}

// HandleGetAll handles GET /api/products.
// Supports query parameters: ?name=searchTerm, ?page=1&limit=20.
func (h *ProductHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	products, err := h.service.GetAll(name)
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}

	page, limit := helper.ParsePagination(r, 20)
	total := len(products)

	// Apply pagination
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
		Items:      products[start:end],
		Page:       page,
		Limit:      limit,
		TotalItems: total,
		TotalPages: totalPages,
	}

	helper.WriteSuccess(w, http.StatusOK, "Success", paged)
}

// HandleGetByID handles GET /api/products/{id}.
func (h *ProductHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/products/", model.ErrProductNotFound)
	if !ok {
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", product)
}

// HandleCreate handles POST /api/products.
func (h *ProductHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var input model.ProductInput
	if !helper.ValidatePayload(w, r, &input) {
		return
	}

	product := &model.Product{
		Name:       input.Name,
		Price:      input.Price,
		Stock:      input.Stock,
		CategoryID: input.CategoryID,
	}
	createdProduct, err := h.service.Create(product)
	if err != nil {
		if errors.Is(err, model.ErrCategoryNotFound) {
			helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Product created successfully", createdProduct)
}

// HandleUpdate handles PUT /api/products/{id}.
func (h *ProductHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/products/", model.ErrProductNotFound)
	if !ok {
		return
	}

	var input model.ProductInput
	if !helper.ValidatePayload(w, r, &input) {
		return
	}

	product := &model.Product{
		Name:       input.Name,
		Price:      input.Price,
		Stock:      input.Stock,
		CategoryID: input.CategoryID,
	}
	updatedProduct, err := h.service.Update(id, product)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		if errors.Is(err, model.ErrCategoryNotFound) {
			helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product updated successfully", updatedProduct)
}

// HandleDelete handles DELETE /api/products/{id}.
func (h *ProductHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/products/", model.ErrProductNotFound)
	if !ok {
		return
	}

	err := h.service.Delete(id)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to delete product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}
