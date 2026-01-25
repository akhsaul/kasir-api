package handler

import (
	"encoding/json"
	"errors"
	"kasir-api/entity"
	"kasir-api/helper"
	"kasir-api/service"
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
func (h *ProductHandler) HandleGetAll(w http.ResponseWriter) {
	products, err := h.service.GetAll()
	if err != nil {
		helper.WriteError(w, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", products)
}

// HandleGetByID handles GET /api/product/{id}
func (h *ProductHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Invalid Product ID", err)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			helper.WriteError(w, http.StatusNotFound, "Product not found", err)
			return
		}
		helper.WriteError(w, http.StatusBadRequest, "Invalid request", err)
		return
	}
	helper.WriteSuccess(w, http.StatusOK, "Success", product)
}

// HandleCreate handles POST /api/product
func (h *ProductHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var product entity.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			// Body close error, already processed the request
		}
	}()

	createdProduct, err := h.service.Create(&product)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Product created successfully", createdProduct)
}

// HandleUpdate handles PUT /api/product/{id}
func (h *ProductHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Invalid Product ID", err)
		return
	}

	var product entity.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Invalid JSON", err)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			// Body close error, already processed the request
		}
	}()

	updatedProduct, err := h.service.Update(id, &product)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			helper.WriteError(w, http.StatusNotFound, "Product not found", err)
			return
		}
		helper.WriteError(w, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product updated successfully", updatedProduct)
}

// HandleDelete handles DELETE /api/product/{id}
func (h *ProductHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/product/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Invalid Product ID", err)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			helper.WriteError(w, http.StatusNotFound, "Product not found", err)
			return
		}
		helper.WriteError(w, http.StatusBadRequest, "Failed to delete product", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}
