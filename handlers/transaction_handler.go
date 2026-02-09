// Package handler: Handler layer — menerima request dan response.
// Error request (invalid ID, invalid JSON, dll) → cek sini.
package handler

import (
	"errors"
	"net/http"
	"time"

	helper "kasir-api/helpers"
	model "kasir-api/models"
	service "kasir-api/services"
)

// TransactionHandler handles HTTP requests for transaction endpoints.
type TransactionHandler struct {
	service *service.TransactionService
}

// NewTransactionHandler creates a new instance of TransactionHandler.
func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: svc,
	}
}

// HandleCheckout handles POST /api/checkout.
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	var request model.CheckoutRequest
	if !helper.ValidatePayload(w, r, &request) {
		return
	}

	transaction, err := h.service.Checkout(&request)
	if err != nil {
		if errors.Is(err, model.ErrProductNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		if errors.Is(err, model.ErrInsufficientStock) ||
			errors.Is(err, model.ErrEmptyCheckout) ||
			errors.Is(err, model.ErrInvalidQuantity) {
			helper.WriteError(w, r, http.StatusBadRequest, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to process checkout", err)
		return
	}

	helper.WriteSuccess(w, http.StatusCreated, "Checkout successful", transaction)
}

// HandleGetByID handles GET /api/transactions/{id}.
func (h *TransactionHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := helper.ParseIDFromPath(w, r, "/api/transactions/", model.ErrTransactionNotFound)
	if !ok {
		return
	}

	transaction, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, model.ErrTransactionNotFound) {
			helper.WriteError(w, r, http.StatusNotFound, err.Error(), err)
			return
		}
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve transaction", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Success", transaction)
}

// HandleGetTodayReport handles GET /api/report/hari-ini.
func (h *TransactionHandler) HandleGetTodayReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve report", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Success", report)
}

// HandleGetReport handles GET /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD.
func (h *TransactionHandler) HandleGetReport(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		helper.WriteError(w, r, http.StatusBadRequest, "start_date and end_date are required", model.ErrInvalidDateRange)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "invalid start_date format, use YYYY-MM-DD", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		helper.WriteError(w, r, http.StatusBadRequest, "invalid end_date format, use YYYY-MM-DD", err)
		return
	}

	if endDate.Before(startDate) {
		helper.WriteError(w, r, http.StatusBadRequest, "end_date must be after start_date", model.ErrInvalidDateRange)
		return
	}

	report, err := h.service.GetReportByDateRange(startDate, endDate)
	if err != nil {
		helper.WriteError(w, r, http.StatusInternalServerError, "Failed to retrieve report", err)
		return
	}

	helper.WriteSuccess(w, http.StatusOK, "Success", report)
}
