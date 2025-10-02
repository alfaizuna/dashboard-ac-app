package handler

import (
	"strconv"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(customerService service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// CreateCustomer creates a new customer
// @Summary Create a new customer
// @Description Create a new customer with the provided information
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body request.CustomerCreateRequest true "Customer creation data"
// @Success 201 {object} response.BaseResponse{data=domain.Customer}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var req request.CustomerCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	customer, err := h.customerService.Create(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create customer")
	}

	return response.Created(c, "Customer created successfully", customer)
}

// GetCustomer retrieves a customer by ID
// @Summary Get customer by ID
// @Description Get a customer by their ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} response.BaseResponse{data=domain.Customer}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Customer ID is required", nil)
	}

	customer, err := h.customerService.GetByID(id)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		return response.InternalServerError(c, "Failed to get customer")
	}

	return response.Success(c, "Customer retrieved successfully", customer)
}

// UpdateCustomer updates an existing customer
// @Summary Update customer
// @Description Update an existing customer with the provided information
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body request.CustomerUpdateRequest true "Customer update data"
// @Success 200 {object} response.BaseResponse{data=domain.Customer}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Customer ID is required", nil)
	}

	var req request.CustomerUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	customer, err := h.customerService.Update(id, &req)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		return response.InternalServerError(c, "Failed to update customer")
	}

	return response.Success(c, "Customer updated successfully", customer)
}

// DeleteCustomer deletes a customer
// @Summary Delete customer
// @Description Delete a customer by their ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Customer ID is required", nil)
	}

	err := h.customerService.Delete(id)
	if err != nil {
		if err.Error() == "customer not found" {
			return response.NotFound(c, "Customer not found")
		}
		return response.InternalServerError(c, "Failed to delete customer")
	}

	return response.Success(c, "Customer deleted successfully", nil)
}

// ListCustomers retrieves a paginated list of customers
// @Summary List customers
// @Description Get a paginated list of customers
// @Tags customers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Customer}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /customers [get]
func (h *CustomerHandler) ListCustomers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	customers, total, err := h.customerService.List(pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get customers")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Customers retrieved successfully", customers, paginationMeta)
}

// SearchCustomers searches for customers based on criteria
// @Summary Search customers
// @Description Search for customers based on name, phone, or email
// @Tags customers
// @Accept json
// @Produce json
// @Param name query string false "Customer name"
// @Param phone query string false "Customer phone"
// @Param email query string false "Customer email"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Customer}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /customers/search [get]
func (h *CustomerHandler) SearchCustomers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	searchReq := &request.CustomerSearchRequest{
		PaginationRequest: &request.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
		Name:  c.Query("name"),
		Phone: c.Query("phone"),
		Email: c.Query("email"),
	}

	customers, total, err := h.customerService.Search(searchReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to search customers")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Customer search completed successfully", customers, paginationMeta)
}