package handler

import (
	"strconv"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type ScheduleHandler struct {
	scheduleService service.ScheduleService
}

func NewScheduleHandler(scheduleService service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// CreateSchedule creates a new schedule
// @Summary Create a new schedule
// @Description Create a new schedule with the provided information
// @Tags schedules
// @Accept json
// @Produce json
// @Param schedule body request.ScheduleCreateRequest true "Schedule creation data"
// @Success 201 {object} response.BaseResponse{data=domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules [post]
func (h *ScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var req request.ScheduleCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	schedule, err := h.scheduleService.Create(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create schedule")
	}

	return response.Created(c, "Schedule created successfully", schedule)
}

// GetSchedule retrieves a schedule by ID
// @Summary Get schedule by ID
// @Description Get a schedule by its ID
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} response.BaseResponse{data=domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/{id} [get]
func (h *ScheduleHandler) GetSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Schedule ID is required", nil)
	}

	schedule, err := h.scheduleService.GetByID(id)
	if err != nil {
		if err.Error() == "schedule not found" {
			return response.NotFound(c, "Schedule not found")
		}
		return response.InternalServerError(c, "Failed to get schedule")
	}

	return response.Success(c, "Schedule retrieved successfully", schedule)
}

// UpdateSchedule updates an existing schedule
// @Summary Update schedule
// @Description Update an existing schedule with the provided information
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Param schedule body request.ScheduleUpdateRequest true "Schedule update data"
// @Success 200 {object} response.BaseResponse{data=domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/{id} [put]
func (h *ScheduleHandler) UpdateSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Schedule ID is required", nil)
	}

	var req request.ScheduleUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	schedule, err := h.scheduleService.Update(id, &req)
	if err != nil {
		if err.Error() == "schedule not found" {
			return response.NotFound(c, "Schedule not found")
		}
		return response.InternalServerError(c, "Failed to update schedule")
	}

	return response.Success(c, "Schedule updated successfully", schedule)
}

// DeleteSchedule deletes a schedule
// @Summary Delete schedule
// @Description Delete a schedule by its ID
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/{id} [delete]
func (h *ScheduleHandler) DeleteSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Schedule ID is required", nil)
	}

	err := h.scheduleService.Delete(id)
	if err != nil {
		if err.Error() == "schedule not found" {
			return response.NotFound(c, "Schedule not found")
		}
		return response.InternalServerError(c, "Failed to delete schedule")
	}

	return response.Success(c, "Schedule deleted successfully", nil)
}

// ListSchedules retrieves a paginated list of schedules
// @Summary List schedules
// @Description Get a paginated list of schedules
// @Tags schedules
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules [get]
func (h *ScheduleHandler) ListSchedules(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	schedules, total, err := h.scheduleService.List(pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get schedules")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Schedules retrieved successfully", schedules, paginationMeta)
}

// SearchSchedules searches for schedules based on criteria
// @Summary Search schedules
// @Description Search for schedules based on customer, technician, service, status, and date range
// @Tags schedules
// @Accept json
// @Produce json
// @Param customer_id query string false "Customer ID"
// @Param technician_id query string false "Technician ID"
// @Param service_id query string false "Service ID"
// @Param status query string false "Schedule status"
// @Param date_from query string false "Date from (YYYY-MM-DD)"
// @Param date_to query string false "Date to (YYYY-MM-DD)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/search [get]
func (h *ScheduleHandler) SearchSchedules(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	searchReq := &request.ScheduleSearchRequest{
		PaginationRequest: &request.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
		CustomerID:   c.Query("customer_id"),
		TechnicianID: c.Query("technician_id"),
		ServiceID:    c.Query("service_id"),
		Status:       c.Query("status"),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
	}

	schedules, total, err := h.scheduleService.Search(searchReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to search schedules")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Schedule search completed successfully", schedules, paginationMeta)
}

// GetSchedulesByCustomer retrieves schedules by customer ID
// @Summary Get schedules by customer
// @Description Get schedules for a specific customer
// @Tags schedules
// @Accept json
// @Produce json
// @Param customer_id path string true "Customer ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/customer/{customer_id} [get]
func (h *ScheduleHandler) GetSchedulesByCustomer(c *fiber.Ctx) error {
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

	schedules, total, err := h.scheduleService.GetByCustomerID(customerID, pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get schedules")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Customer schedules retrieved successfully", schedules, paginationMeta)
}

// GetSchedulesByTechnician retrieves schedules by technician ID
// @Summary Get schedules by technician
// @Description Get schedules for a specific technician
// @Tags schedules
// @Accept json
// @Produce json
// @Param technician_id path string true "Technician ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/technician/{technician_id} [get]
func (h *ScheduleHandler) GetSchedulesByTechnician(c *fiber.Ctx) error {
	technicianID := c.Params("technician_id")
	if technicianID == "" {
		return response.BadRequest(c, "Technician ID is required", nil)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	schedules, total, err := h.scheduleService.GetByTechnicianID(technicianID, pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get schedules")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Technician schedules retrieved successfully", schedules, paginationMeta)
}

// GetSchedulesByStatus retrieves schedules by status
// @Summary Get schedules by status
// @Description Get schedules with a specific status
// @Tags schedules
// @Accept json
// @Produce json
// @Param status path string true "Schedule status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Schedule}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /schedules/status/{status} [get]
func (h *ScheduleHandler) GetSchedulesByStatus(c *fiber.Ctx) error {
	statusStr := c.Params("status")
	if statusStr == "" {
		return response.BadRequest(c, "Status is required", nil)
	}

	status := domain.ScheduleStatus(statusStr)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	schedules, total, err := h.scheduleService.GetByStatus(status, pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get schedules")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Schedules by status retrieved successfully", schedules, paginationMeta)
}