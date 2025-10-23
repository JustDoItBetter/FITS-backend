package signing

import (
	"github.com/JustDoItBetter/FITS-backend/internal/common/response"
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for signing endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new signing handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes registers signing routes
func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/upload", h.Upload)
	router.Get("/sign_requests", h.GetSignRequests)
	router.Post("/sign_uploads", h.SignUploads)
}

// Upload handles parquet file uploads
// @Summary Upload parquet file
// @Description Upload a parquet file containing student data
// @Tags signing
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Parquet file"
// @Success 201 {object} response.SuccessResponse{data=UploadRecord}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 501 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/signing/upload [post]
func (h *Handler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, err)
	}

	// Read file data
	fileData, err := file.Open()
	if err != nil {
		return response.Error(c, err)
	}
	defer fileData.Close()

	data := make([]byte, file.Size)
	if _, err := fileData.Read(data); err != nil {
		return response.Error(c, err)
	}

	record, err := h.service.HandleUpload(c.Context(), file.Filename, data)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, record)
}

// GetSignRequests retrieves pending sign requests
// @Summary Get pending sign requests
// @Description Retrieve all pending sign requests as a parquet file
// @Tags signing
// @Produce application/octet-stream
// @Success 200 {file} binary "Parquet file with pending sign requests"
// @Failure 401 {object} response.ErrorResponse
// @Failure 501 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/signing/sign_requests [get]
func (h *Handler) GetSignRequests(c *fiber.Ctx) error {
	data, err := h.service.GetPendingSignRequests(c.Context())
	if err != nil {
		return response.Error(c, err)
	}

	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename=sign_requests.parquet")

	return c.Send(data)
}

// SignUploads handles signed request uploads
// @Summary Upload signed requests
// @Description Upload a parquet file containing signed requests
// @Tags signing
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Parquet file with signed requests"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 501 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/signing/sign_uploads [post]
func (h *Handler) SignUploads(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Error(c, err)
	}

	// Read file data
	fileData, err := file.Open()
	if err != nil {
		return response.Error(c, err)
	}
	defer fileData.Close()

	data := make([]byte, file.Size)
	if _, err := fileData.Read(data); err != nil {
		return response.Error(c, err)
	}

	if err := h.service.HandleSignedUploads(c.Context(), file.Filename, data); err != nil {
		return response.Error(c, err)
	}

	return response.SuccessWithMessage(c, "signed requests processed successfully", nil)
}
