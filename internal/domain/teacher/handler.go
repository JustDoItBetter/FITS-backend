package teacher

import (
	"github.com/JustDoItBetter/FITS-backend/internal/common/pagination"
	"github.com/JustDoItBetter/FITS-backend/internal/common/response"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/", h.Create)
	router.Get("/:uuid", h.GetByUUID)
	router.Put("/:uuid", h.Update)
	router.Delete("/:uuid", h.Delete)
	router.Get("/", h.List)
}

// Create godoc
// @Summary Create a new teacher
// @Description Creates a new teacher record. Requires admin role. Email must be unique. Department is required.
// @Tags Teachers
// @Accept json
// @Produce json
// @Param request body CreateTeacherRequest true "Teacher creation request"
// @Success 201 {object} response.SuccessResponse{data=Teacher} "Teacher created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - requires admin role"
// @Failure 409 {object} response.ErrorResponse "Conflict - email already exists"
// @Failure 422 {object} response.ErrorResponse "Validation error - invalid field values"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/teacher [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateTeacherRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	teacher, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, teacher)
}

// GetByUUID godoc
// @Summary Get teacher by UUID
// @Description Retrieves detailed information about a specific teacher by their UUID. Public endpoint.
// @Tags Teachers
// @Produce json
// @Param uuid path string true "Teacher UUID" format(uuid) example(550e8400-e29b-41d4-a716-446655440010)
// @Success 200 {object} response.SuccessResponse{data=Teacher} "Teacher found"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 404 {object} response.ErrorResponse "Teacher not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/teacher/{uuid} [get]
func (h *Handler) GetByUUID(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	teacher, err := h.service.GetByUUID(c.Context(), uuid)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, teacher)
}

// Update godoc
// @Summary Update teacher information
// @Description Updates an existing teacher's information. Requires admin role. Supports partial updates including department changes.
// @Tags Teachers
// @Accept json
// @Produce json
// @Param uuid path string true "Teacher UUID" format(uuid) example(550e8400-e29b-41d4-a716-446655440010)
// @Param request body UpdateTeacherRequest true "Teacher update request (partial updates supported)"
// @Success 200 {object} response.SuccessResponse{data=Teacher} "Teacher updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or UUID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - requires admin role"
// @Failure 404 {object} response.ErrorResponse "Teacher not found"
// @Failure 409 {object} response.ErrorResponse "Conflict - email already exists"
// @Failure 422 {object} response.ErrorResponse "Validation error - invalid field values"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/teacher/{uuid} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	var req UpdateTeacherRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	teacher, err := h.service.Update(c.Context(), uuid, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, teacher)
}

// Delete godoc
// @Summary Delete a teacher
// @Description Permanently deletes a teacher from the system (soft delete). Requires admin role.
// @Tags Teachers
// @Produce json
// @Param uuid path string true "Teacher UUID" format(uuid) example(550e8400-e29b-41d4-a716-446655440010)
// @Success 204 "Teacher deleted successfully (no content)"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - requires admin role"
// @Failure 404 {object} response.ErrorResponse "Teacher not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/teacher/{uuid} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	if err := h.service.Delete(c.Context(), uuid); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// List godoc
// @Summary List teachers with pagination
// @Description Retrieves a paginated list of teachers. Supports page and limit query parameters.
// @Description Default: page=1, limit=20. Maximum limit is 100 to prevent performance issues.
// @Tags Teachers
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 20, max: 100)" minimum(1) maximum(100)
// @Success 200 {object} pagination.Response{data=[]Teacher} "Paginated list of teachers with metadata"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/teacher [get]
func (h *Handler) List(c *fiber.Ctx) error {
	// Extract pagination parameters from query string
	params := pagination.ExtractParams(c)

	teachers, totalCount, err := h.service.ListPaginated(c.Context(), params)
	if err != nil {
		return response.Error(c, err)
	}

	// Build paginated response with metadata
	paginatedResp := pagination.NewResponse(teachers, params, totalCount)
	return c.JSON(paginatedResp)
}
