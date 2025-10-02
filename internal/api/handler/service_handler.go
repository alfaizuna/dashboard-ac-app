package handler

import (
	"strconv"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ServiceHandler struct {
	serviceService service.ServiceService
}

func NewServiceHandler(serviceService service.ServiceService) *ServiceHandler {
	return &ServiceHandler{
		serviceService: serviceService,
	}
}

// CreateService creates a new service
// @Summary Create a new service
// @Description Create a new service with the provided information
// @Tags services
// @Accept json
// @Produce json
// @Param service body request.ServiceCreateRequest true "Service creation data"
// @Success 201 {object} response.BaseResponse{data=domain.Service}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /services [post]
func (h *ServiceHandler) CreateService(c *fiber.Ctx) error {
	var req request.ServiceCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	service, err := h.serviceService.Create(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create service")
	}

	return response.Created(c, "Service created successfully", service)
}

// GetService retrieves a service by ID
// @Summary Get service by ID
// @Description Get a service by its ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} response.BaseResponse{data=domain.Service}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /services/{id} [get]
func (h *ServiceHandler) GetService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Service ID is required", nil)
	}

	service, err := h.serviceService.GetByID(id)
	if err != nil {
		if err.Error() == "service not found" {
			return response.NotFound(c, "Service not found")
		}
		return response.InternalServerError(c, "Failed to get service")
	}

	return response.Success(c, "Service retrieved successfully", service)
}

// UpdateService updates an existing service
// @Summary Update service
// @Description Update an existing service with the provided information
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param service body request.ServiceUpdateRequest true "Service update data"
// @Success 200 {object} response.BaseResponse{data=domain.Service}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /services/{id} [put]
func (h *ServiceHandler) UpdateService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Service ID is required", nil)
	}

	var req request.ServiceUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	service, err := h.serviceService.Update(id, &req)
	if err != nil {
		if err.Error() == "service not found" {
			return response.NotFound(c, "Service not found")
		}
		return response.InternalServerError(c, "Failed to update service")
	}

	return response.Success(c, "Service updated successfully", service)
}

// DeleteService deletes a service
// @Summary Delete service
// @Description Delete a service by its ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /services/{id} [delete]
func (h *ServiceHandler) DeleteService(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Service ID is required", nil)
	}

	err := h.serviceService.Delete(id)
	if err != nil {
		if err.Error() == "service not found" {
			return response.NotFound(c, "Service not found")
		}
		return response.InternalServerError(c, "Failed to delete service")
	}

	return response.Success(c, "Service deleted successfully", nil)
}

// ListServices retrieves a paginated list of services
// @Summary List services
// @Description Get a paginated list of services
// @Tags services
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Service}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /services [get]
func (h *ServiceHandler) ListServices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	services, total, err := h.serviceService.List(pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get services")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Services retrieved successfully", services, paginationMeta)
}

// SearchServices searches for services based on criteria
// @Summary Search services
// @Description Search for services based on name and price range
// @Tags services
// @Accept json
// @Produce json
// @Param name query string false "Service name"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Service}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /services/search [get]
func (h *ServiceHandler) SearchServices(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	
	minPrice, _ := strconv.ParseFloat(c.Query("min_price", "0"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price", "0"), 64)

	searchReq := &request.ServiceSearchRequest{
		PaginationRequest: &request.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
		Name:     c.Query("name"),
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}

	services, total, err := h.serviceService.Search(searchReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to search services")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Service search completed successfully", services, paginationMeta)
}