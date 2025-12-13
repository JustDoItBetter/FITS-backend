// Package student provides handlers and services for student management.
package student

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

// RegisterRoutes registers all student endpoints with their required middleware
// This provides a single source of truth for routes and their security requirements
func (h *Handler) RegisterRoutes(router fiber.Router, jwtMW JWTMiddleware, rbacMW RBACMiddleware) {
	// POST /api/v1/student - Create student (Admin only)
	router.Post("/",
		jwtMW.RequireAuth(),
		rbacMW.RequireAdmin(),
		h.Create,
	)

	// GET /api/v1/student/:uuid - Get student (optional auth for future filtering)
	router.Get("/:uuid",
		jwtMW.OptionalAuth(),
		h.GetByUUID,
	)

	// PUT /api/v1/student/:uuid - Update student (Admin only)
	router.Put("/:uuid",
		jwtMW.RequireAuth(),
		rbacMW.RequireAdmin(),
		h.Update,
	)

	// DELETE /api/v1/student/:uuid - Delete student (Admin only, soft delete)
	router.Delete("/:uuid",
		jwtMW.RequireAuth(),
		rbacMW.RequireAdmin(),
		h.Delete,
	)

	// GET /api/v1/student - List students (optional auth for future filtering)
	router.Get("/",
		jwtMW.OptionalAuth(),
		h.List,
	)
}

// JWTMiddleware interface defines JWT authentication middleware requirements
type JWTMiddleware interface {
	RequireAuth() fiber.Handler
	OptionalAuth() fiber.Handler
}

// RBACMiddleware interface defines role-based access control middleware
type RBACMiddleware interface {
	RequireAdmin() fiber.Handler
	RequireRole(roles ...interface{}) fiber.Handler
}

// Create godoc
// @Summary Create a new student
// @Description Creates a new student record. Requires admin role. Email must be unique.
// @Tags Students
// @Accept json
// @Produce json
// @Param request body CreateStudentRequest true "Student creation request"
// @Success 201 {object} response.SuccessResponse{data=Student} "Student created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - requires admin role"
// @Failure 409 {object} response.ErrorResponse "Conflict - email already exists"
// @Failure 422 {object} response.ErrorResponse "Validation error - invalid field values"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/student [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	student, err := h.service.Create(c.Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, student)
}

// GetByUUID godoc
// @Summary Get student by UUID
// @Description Retrieves detailed information about a specific student by their UUID. Public endpoint.
// @Tags Students
// @Produce json
// @Param uuid path string true "Student UUID" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 200 {object} response.SuccessResponse{data=Student} "Student found"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 404 {object} response.ErrorResponse "Student not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/student/{uuid} [get]
func (h *Handler) GetByUUID(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	student, err := h.service.GetByUUID(c.Context(), uuid)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, student)
}

// Update godoc
// @Summary Update student information
// @Description Updates an existing student's information. Requires admin role. Supports partial updates.
// @Tags Students
// @Accept json
// @Produce json
// @Param uuid path string true "Student UUID" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Param request body UpdateStudentRequest true "Student update request (partial updates supported)"
// @Success 200 {object} response.SuccessResponse{data=Student} "Student updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or UUID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - requires admin role"
// @Failure 404 {object} response.ErrorResponse "Student not found"
// @Failure 409 {object} response.ErrorResponse "Conflict - email already exists"
// @Failure 422 {object} response.ErrorResponse "Validation error - invalid field values"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/student/{uuid} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	var req UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	student, err := h.service.Update(c.Context(), uuid, &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, student)
}

// Delete godoc
// @Summary Delete a student
// @Description Permanently deletes a student from the system (soft delete). Requires admin role.
// @Tags Students
// @Produce json
// @Param uuid path string true "Student UUID" format(uuid) example(550e8400-e29b-41d4-a716-446655440000)
// @Success 204 "Student deleted successfully (no content)"
// @Failure 400 {object} response.ErrorResponse "Invalid UUID format"
// @Failure 401 {object} response.ErrorResponse "Unauthorized - missing or invalid token"
// @Failure 403 {object} response.ErrorResponse "Forbidden - requires admin role"
// @Failure 404 {object} response.ErrorResponse "Student not found"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/student/{uuid} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	if err := h.service.Delete(c.Context(), uuid); err != nil {
		return response.Error(c, err)
	}

	return response.NoContent(c)
}

// List godoc
// @Summary List students with pagination
// @Description Retrieves a paginated list of students. Supports page and limit query parameters.
// @Description Default: page=1, limit=20. Maximum limit is 100 to prevent performance issues.
// @Tags Students
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param limit query int false "Items per page (default: 20, max: 100)" minimum(1) maximum(100)
// @Success 200 {object} pagination.Response{data=[]Student} "Paginated list of students with metadata"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/student [get]
func (h *Handler) List(c *fiber.Ctx) error {
	// Extract pagination parameters from query string
	params := pagination.ExtractParams(c)

	students, totalCount, err := h.service.ListPaginated(c.Context(), params)
	if err != nil {
		return response.Error(c, err)
	}

	// Build paginated response with metadata
	paginatedResp := pagination.NewResponse(students, params, totalCount)
	return c.JSON(paginatedResp)
}
