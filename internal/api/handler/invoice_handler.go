package handler

import (
	"strconv"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type InvoiceHandler struct {
	invoiceService service.InvoiceService
}

func NewInvoiceHandler(invoiceService service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{
		invoiceService: invoiceService,
	}
}

// CreateInvoice creates a new invoice
// @Summary Create a new invoice
// @Description Create a new invoice with the provided information
// @Tags invoices
// @Accept json
// @Produce json
// @Param invoice body request.InvoiceCreateRequest true "Invoice creation data"
// @Success 201 {object} response.BaseResponse{data=domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices [post]
func (h *InvoiceHandler) CreateInvoice(c *fiber.Ctx) error {
	var req request.InvoiceCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	invoice, err := h.invoiceService.Create(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create invoice")
	}

	return response.Created(c, "Invoice created successfully", invoice)
}

// GetInvoice retrieves an invoice by ID
// @Summary Get invoice by ID
// @Description Get an invoice by its ID
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} response.BaseResponse{data=domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/{id} [get]
func (h *InvoiceHandler) GetInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Invoice ID is required", nil)
	}

	invoice, err := h.invoiceService.GetByID(id)
	if err != nil {
		if err.Error() == "invoice not found" {
			return response.NotFound(c, "Invoice not found")
		}
		return response.InternalServerError(c, "Failed to get invoice")
	}

	return response.Success(c, "Invoice retrieved successfully", invoice)
}

// UpdateInvoice updates an existing invoice
// @Summary Update invoice
// @Description Update an existing invoice with the provided information
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Param invoice body request.InvoiceUpdateRequest true "Invoice update data"
// @Success 200 {object} response.BaseResponse{data=domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/{id} [put]
func (h *InvoiceHandler) UpdateInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Invoice ID is required", nil)
	}

	var req request.InvoiceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	invoice, err := h.invoiceService.Update(id, &req)
	if err != nil {
		if err.Error() == "invoice not found" {
			return response.NotFound(c, "Invoice not found")
		}
		return response.InternalServerError(c, "Failed to update invoice")
	}

	return response.Success(c, "Invoice updated successfully", invoice)
}

// DeleteInvoice deletes an invoice
// @Summary Delete invoice
// @Description Delete an invoice by its ID
// @Tags invoices
// @Accept json
// @Produce json
// @Param id path string true "Invoice ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/{id} [delete]
func (h *InvoiceHandler) DeleteInvoice(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Invoice ID is required", nil)
	}

	err := h.invoiceService.Delete(id)
	if err != nil {
		if err.Error() == "invoice not found" {
			return response.NotFound(c, "Invoice not found")
		}
		return response.InternalServerError(c, "Failed to delete invoice")
	}

	return response.Success(c, "Invoice deleted successfully", nil)
}

// ListInvoices retrieves a paginated list of invoices
// @Summary List invoices
// @Description Get a paginated list of invoices
// @Tags invoices
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices [get]
func (h *InvoiceHandler) ListInvoices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	invoices, total, err := h.invoiceService.List(pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get invoices")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Invoices retrieved successfully", invoices, paginationMeta)
}

// SearchInvoices searches for invoices based on criteria
// @Summary Search invoices
// @Description Search for invoices based on customer, schedule, status, and date range
// @Tags invoices
// @Accept json
// @Produce json
// @Param customer_id query string false "Customer ID"
// @Param schedule_id query string false "Schedule ID"
// @Param status query string false "Invoice status"
// @Param date_from query string false "Date from (YYYY-MM-DD)"
// @Param date_to query string false "Date to (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/search [get]
func (h *InvoiceHandler) SearchInvoices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	searchReq := &request.InvoiceSearchRequest{
		PaginationRequest: &request.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
		CustomerID: c.Query("customer_id"),
		ScheduleID: c.Query("schedule_id"),
		Status:     c.Query("status"),
		DateFrom:   c.Query("date_from"),
		DateTo:     c.Query("date_to"),
	}

	invoices, total, err := h.invoiceService.Search(searchReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to search invoices")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Invoice search completed successfully", invoices, paginationMeta)
}

// GetInvoicesByCustomer retrieves invoices by customer ID
// @Summary Get invoices by customer
// @Description Get invoices for a specific customer
// @Tags invoices
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/customer/{customer_id} [get]
func (h *InvoiceHandler) GetInvoicesByCustomer(c *fiber.Ctx) error {
	customerID := c.Params("customer_id")
	if customerID == "" {
		return response.BadRequest(c, "Customer ID is required", nil)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	invoices, total, err := h.invoiceService.GetByCustomerID(customerID, pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get invoices")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Customer invoices retrieved successfully", invoices, paginationMeta)
}

// GetInvoicesBySchedule retrieves invoices by schedule ID
// @Summary Get invoices by schedule
// @Description Get invoices for a specific schedule
// @Tags invoices
// @Accept json
// @Produce json
// @Param schedule_id path string true "Schedule ID"
// @Success 200 {object} response.BaseResponse{data=domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/schedule/{schedule_id} [get]
func (h *InvoiceHandler) GetInvoicesBySchedule(c *fiber.Ctx) error {
	scheduleID := c.Params("schedule_id")
	if scheduleID == "" {
		return response.BadRequest(c, "Schedule ID is required", nil)
	}

	invoice, err := h.invoiceService.GetByScheduleID(scheduleID)
	if err != nil {
		if err.Error() == "invoice not found" {
			return response.NotFound(c, "Invoice not found")
		}
		return response.InternalServerError(c, "Failed to get invoice")
	}

	return response.Success(c, "Schedule invoice retrieved successfully", invoice)
}

// GetInvoicesByStatus retrieves invoices by status
// @Summary Get invoices by status
// @Description Get invoices with a specific status
// @Tags invoices
// @Accept json
// @Produce json
// @Param status path string true "Invoice status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Invoice}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /invoices/status/{status} [get]
func (h *InvoiceHandler) GetInvoicesByStatus(c *fiber.Ctx) error {
	statusStr := c.Params("status")
	if statusStr == "" {
		return response.BadRequest(c, "Status is required", nil)
	}

	status := domain.InvoiceStatus(statusStr)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	invoices, total, err := h.invoiceService.GetByStatus(status, pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get invoices")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Invoices by status retrieved successfully", invoices, paginationMeta)
}