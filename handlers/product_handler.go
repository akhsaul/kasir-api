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

// HandleGetAll handles GET /api/product
func (h *ProductHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", products)
}

// HandleGetByID handles GET /api/product/{id}
func (h *ProductHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid Product ID", err)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid request", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", product)
}

// HandleCreate handles POST /api/product
func (h *ProductHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if !helper.ValidatePayload(w, r, &product) {
		return
	}

	createdProduct, err := h.service.Create(&product)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Product created successfully", createdProduct)
}

// HandleUpdate handles PUT /api/product/{id}
func (h *ProductHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid Product ID", err)
		return
	}

	var product model.Product
	if !helper.ValidatePayload(w, r, &product) {
		return
	}

	updatedProduct, err := h.service.Update(id, &product)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product updated successfully", updatedProduct)
}

// HandleDelete handles DELETE /api/product/{id}
func (h *ProductHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "Invalid Product ID", err)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, "Product not found", err)
			return
		}
		helper.WriteError(w, r, http.StatusBadRequest, "Failed to delete product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}
