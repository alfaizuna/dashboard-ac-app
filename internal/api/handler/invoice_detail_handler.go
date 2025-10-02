package handler

import (
	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type InvoiceDetailHandler struct {
	invoiceDetailService service.InvoiceDetailService
}

func NewInvoiceDetailHandler(invoiceDetailService service.InvoiceDetailService) *InvoiceDetailHandler {
	return &InvoiceDetailHandler{
		invoiceDetailService: invoiceDetailService,
	}
}

// CreateInvoiceDetail creates a new invoice detail
// @Summary Create a new invoice detail
// @Description Create a new invoice detail with the provided information
// @Tags invoice-details
// @Accept json
// @Produce json
// @Param invoice_detail body request.InvoiceDetailCreateRequest true "Invoice detail creation data"
// @Success 201 {object} response.BaseResponse{data=domain.InvoiceDetail}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoice-details [post]
func (h *InvoiceDetailHandler) CreateInvoiceDetail(c *fiber.Ctx) error {
	var req request.InvoiceDetailCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	invoiceDetail, err := h.invoiceDetailService.Create(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create invoice detail")
	}

	return response.Created(c, "Invoice detail created successfully", invoiceDetail)
}

// GetInvoiceDetail retrieves an invoice detail by ID
// @Summary Get invoice detail by ID
// @Description Get an invoice detail by its ID
// @Tags invoice-details
// @Accept json
// @Produce json
// @Param id path string true "Invoice detail ID"
// @Success 200 {object} response.BaseResponse{data=domain.InvoiceDetail}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoice-details/{id} [get]
func (h *InvoiceDetailHandler) GetInvoiceDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Invoice detail ID is required", nil)
	}

	invoiceDetail, err := h.invoiceDetailService.GetByID(id)
	if err != nil {
		if err.Error() == "invoice detail not found" {
			return response.NotFound(c, "Invoice detail not found")
		}
		return response.InternalServerError(c, "Failed to get invoice detail")
	}

	return response.Success(c, "Invoice detail retrieved successfully", invoiceDetail)
}

// UpdateInvoiceDetail updates an existing invoice detail
// @Summary Update invoice detail
// @Description Update an existing invoice detail with the provided information
// @Tags invoice-details
// @Accept json
// @Produce json
// @Param id path string true "Invoice detail ID"
// @Param invoice_detail body request.InvoiceDetailUpdateRequest true "Invoice detail update data"
// @Success 200 {object} response.BaseResponse{data=domain.InvoiceDetail}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoice-details/{id} [put]
func (h *InvoiceDetailHandler) UpdateInvoiceDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Invoice detail ID is required", nil)
	}

	var req request.InvoiceDetailUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	invoiceDetail, err := h.invoiceDetailService.Update(id, &req)
	if err != nil {
		if err.Error() == "invoice detail not found" {
			return response.NotFound(c, "Invoice detail not found")
		}
		return response.InternalServerError(c, "Failed to update invoice detail")
	}

	return response.Success(c, "Invoice detail updated successfully", invoiceDetail)
}

// DeleteInvoiceDetail deletes an invoice detail
// @Summary Delete invoice detail
// @Description Delete an invoice detail by its ID
// @Tags invoice-details
// @Accept json
// @Produce json
// @Param id path string true "Invoice detail ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoice-details/{id} [delete]
func (h *InvoiceDetailHandler) DeleteInvoiceDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Invoice detail ID is required", nil)
	}

	err := h.invoiceDetailService.Delete(id)
	if err != nil {
		if err.Error() == "invoice detail not found" {
			return response.NotFound(c, "Invoice detail not found")
		}
		return response.InternalServerError(c, "Failed to delete invoice detail")
	}

	return response.Success(c, "Invoice detail deleted successfully", nil)
}

// GetInvoiceDetailsByInvoice retrieves invoice details by invoice ID
// @Summary Get invoice details by invoice
// @Description Get all invoice details for a specific invoice
// @Tags invoice-details
// @Accept json
// @Produce json
// @Param invoice_id path string true "Invoice ID"
// @Success 200 {object} response.BaseResponse{data=[]domain.InvoiceDetail}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoice-details/invoice/{invoice_id} [get]
func (h *InvoiceDetailHandler) GetInvoiceDetailsByInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("invoice_id")
	if invoiceID == "" {
		return response.BadRequest(c, "Invoice ID is required", nil)
	}

	invoiceDetails, err := h.invoiceDetailService.GetByInvoiceID(invoiceID)
	if err != nil {
		return response.InternalServerError(c, "Failed to get invoice details")
	}

	return response.Success(c, "Invoice details retrieved successfully", invoiceDetails)
}

// DeleteInvoiceDetailsByInvoice deletes all invoice details for an invoice
// @Summary Delete invoice details by invoice
// @Description Delete all invoice details for a specific invoice
// @Tags invoice-details
// @Accept json
// @Produce json
// @Param invoice_id path string true "Invoice ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoice-details/invoice/{invoice_id} [delete]
func (h *InvoiceDetailHandler) DeleteInvoiceDetailsByInvoice(c *fiber.Ctx) error {
	invoiceID := c.Params("invoice_id")
	if invoiceID == "" {
		return response.BadRequest(c, "Invoice ID is required", nil)
	}

	err := h.invoiceDetailService.DeleteByInvoiceID(invoiceID)
	if err != nil {
		return response.InternalServerError(c, "Failed to delete invoice details")
	}

	return response.Success(c, "Invoice details deleted successfully", nil)
}