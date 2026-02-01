// Package handler: Handler layer — menerima request dan response.
// Error request (invalid ID, invalid JSON, dll) → cek sini.
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

// ProductHandler handles HTTP requests for product endpoints
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler creates a new instance of ProductHandler
func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// HandleGetAll handles GET /api/products
func (h *ProductHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", products)
}

// HandleGetByID handles GET /api/products/{id}
func (h *ProductHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusNotFound, model.ErrProductNotFound.Error(), model.ErrProductNotFound)
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

// HandleCreate handles POST /api/products
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
		if errors.Is(err, model.ErrPriceInvalid) || errors.Is(err, model.ErrStockInvalid) || errors.Is(err, model.ErrNameRequired) {
			helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Product created successfully", createdProduct)
}

// HandleUpdate handles PUT /api/products/{id}
func (h *ProductHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusNotFound, model.ErrProductNotFound.Error(), model.ErrProductNotFound)
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
		if errors.Is(err, model.ErrPriceInvalid) || errors.Is(err, model.ErrStockInvalid) || errors.Is(err, model.ErrNameRequired) {
			helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product updated successfully", updatedProduct)
}

// HandleDelete handles DELETE /api/products/{id}
func (h *ProductHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusNotFound, model.ErrProductNotFound.Error(), model.ErrProductNotFound)
		return
	}

	err = h.service.Delete(id)
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
