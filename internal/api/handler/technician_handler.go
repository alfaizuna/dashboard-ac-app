package handler

import (
	"strconv"

	"dashboard-ac-backend/internal/api/request"
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

type TechnicianHandler struct {
	technicianService service.TechnicianService
}

func NewTechnicianHandler(technicianService service.TechnicianService) *TechnicianHandler {
	return &TechnicianHandler{
		technicianService: technicianService,
	}
}

// CreateTechnician creates a new technician
// @Summary Create a new technician
// @Description Create a new technician with the provided information
// @Tags technicians
// @Accept json
// @Produce json
// @Param technician body request.TechnicianCreateRequest true "Technician creation data"
// @Success 201 {object} response.BaseResponse{data=domain.Technician}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /technicians [post]
func (h *TechnicianHandler) CreateTechnician(c *fiber.Ctx) error {
	var req request.TechnicianCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	technician, err := h.technicianService.Create(&req)
	if err != nil {
		return response.InternalServerError(c, "Failed to create technician")
	}

	return response.Created(c, "Technician created successfully", technician)
}

// GetTechnician retrieves a technician by ID
// @Summary Get technician by ID
// @Description Get a technician by their ID
// @Tags technicians
// @Accept json
// @Produce json
// @Param id path string true "Technician ID"
// @Success 200 {object} response.BaseResponse{data=domain.Technician}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /technicians/{id} [get]
func (h *TechnicianHandler) GetTechnician(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Technician ID is required", nil)
	}

	technician, err := h.technicianService.GetByID(id)
	if err != nil {
		if err.Error() == "technician not found" {
			return response.NotFound(c, "Technician not found")
		}
		return response.InternalServerError(c, "Failed to get technician")
	}

	return response.Success(c, "Technician retrieved successfully", technician)
}

// UpdateTechnician updates an existing technician
// @Summary Update technician
// @Description Update an existing technician with the provided information
// @Tags technicians
// @Accept json
// @Produce json
// @Param id path string true "Technician ID"
// @Param technician body request.TechnicianUpdateRequest true "Technician update data"
// @Success 200 {object} response.BaseResponse{data=domain.Technician}
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /technicians/{id} [put]
func (h *TechnicianHandler) UpdateTechnician(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Technician ID is required", nil)
	}

	var req request.TechnicianUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	technician, err := h.technicianService.Update(id, &req)
	if err != nil {
		if err.Error() == "technician not found" {
			return response.NotFound(c, "Technician not found")
		}
		return response.InternalServerError(c, "Failed to update technician")
	}

	return response.Success(c, "Technician updated successfully", technician)
}

// DeleteTechnician deletes a technician
// @Summary Delete technician
// @Description Delete a technician by their ID
// @Tags technicians
// @Accept json
// @Produce json
// @Param id path string true "Technician ID"
// @Success 200 {object} response.BaseResponse
// @Failure 400 {object} response.BaseResponse
// @Failure 404 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /technicians/{id} [delete]
func (h *TechnicianHandler) DeleteTechnician(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.BadRequest(c, "Technician ID is required", nil)
	}

	err := h.technicianService.Delete(id)
	if err != nil {
		if err.Error() == "technician not found" {
			return response.NotFound(c, "Technician not found")
		}
		return response.InternalServerError(c, "Failed to delete technician")
	}

	return response.Success(c, "Technician deleted successfully", nil)
}

// ListTechnicians retrieves a paginated list of technicians
// @Summary List technicians
// @Description Get a paginated list of technicians
// @Tags technicians
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Technician}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /technicians [get]
func (h *TechnicianHandler) ListTechnicians(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pagination := &request.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	technicians, total, err := h.technicianService.List(pagination)
	if err != nil {
		return response.InternalServerError(c, "Failed to get technicians")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Technicians retrieved successfully", technicians, paginationMeta)
}

// SearchTechnicians searches for technicians based on criteria
// @Summary Search technicians
// @Description Search for technicians based on name or specialization
// @Tags technicians
// @Accept json
// @Produce json
// @Param name query string false "Technician name"
// @Param specialization query string false "Technician specialization"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.BaseResponse{data=[]domain.Technician}
// @Failure 400 {object} response.BaseResponse
// @Failure 500 {object} response.BaseResponse
// @Router /technicians/search [get]
func (h *TechnicianHandler) SearchTechnicians(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	searchReq := &request.TechnicianSearchRequest{
		PaginationRequest: &request.PaginationRequest{
			Page:  page,
			Limit: limit,
		},
		Name:           c.Query("name"),
		Specialization: c.Query("specialization"),
	}

	technicians, total, err := h.technicianService.Search(searchReq)
	if err != nil {
		return response.InternalServerError(c, "Failed to search technicians")
	}

	paginationMeta := response.CalculatePagination(page, limit, total)
	return response.Paginated(c, "Technician search completed successfully", technicians, paginationMeta)
}