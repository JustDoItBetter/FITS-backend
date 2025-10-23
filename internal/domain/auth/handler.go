package auth

import (
	"github.com/JustDoItBetter/FITS-backend/internal/common/response"
	"github.com/gofiber/fiber/v2"
)

// Handler handles HTTP requests for auth endpoints
type Handler struct {
	bootstrapService  *BootstrapService
	invitationService *InvitationService
	authService       *AuthService
}

// NewHandler creates a new auth handler
func NewHandler(bootstrapService *BootstrapService, invitationService *InvitationService, authService *AuthService) *Handler {
	return &Handler{
		bootstrapService:  bootstrapService,
		invitationService: invitationService,
		authService:       authService,
	}
}

// RegisterRoutes registers all auth-related routes
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Bootstrap routes (no auth required)
	bootstrap := app.Group("/api/v1/bootstrap")
	bootstrap.Post("/init", h.InitializeAdmin)

	// Auth routes (no auth required for login/refresh)
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)
	// Logout requires auth - will be added with middleware

	// Invitation routes (public for getting details and completing)
	invite := app.Group("/api/v1/invite")
	invite.Get("/:token", h.GetInvitationDetails)
	invite.Post("/:token/complete", h.CompleteInvitation)

	// Note: Admin invitation creation (/api/v1/admin/invite) is registered
	// in main.go with RequireAuth + RequireAdmin middleware
}

// Bootstrap Endpoints

// InitializeAdmin handles admin initialization
// @Summary Initialize admin certificate
// @Description Generates admin RSA keypair and returns admin token. Can only be called once.
// @Tags bootstrap
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=BootstrapResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/bootstrap/init [post]
func (h *Handler) InitializeAdmin(c *fiber.Ctx) error {
	result, err := h.bootstrapService.InitializeAdmin(c.Context())
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result)
}

// Auth Endpoints

// Login handles user authentication
// @Summary User login
// @Description Authenticate with username and password, returns access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} response.SuccessResponse{data=LoginResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	result, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result)
}

// RefreshToken handles access token refresh
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.SuccessResponse{data=LoginResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	result, err := h.authService.RefreshAccessToken(c.Context(), req.RefreshToken)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result)
}

// Logout handles user logout
// @Summary User logout
// @Description Logout user by invalidating all refresh tokens
// @Tags auth
// @Produce json
// @Success 200 {object} response.SuccessResponse
// @Failure 401 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	// Get user ID from context (set by JWT middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		return response.Error(c, fiber.NewError(fiber.StatusUnauthorized, "user not authenticated"))
	}

	if err := h.authService.Logout(c.Context(), userID.(string)); err != nil {
		return response.Error(c, err)
	}

	return response.SuccessWithMessage(c, "logged out successfully", nil)
}

// Invitation Endpoints

// CreateInvitation creates a new user invitation
// @Summary Create invitation
// @Description Create invitation link for student or teacher (Admin only)
// @Tags invitations
// @Accept json
// @Produce json
// @Param invitation body CreateInvitationRequest true "Invitation data"
// @Success 201 {object} response.SuccessResponse{data=CreateInvitationResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /api/v1/admin/invite [post]
func (h *Handler) CreateInvitation(c *fiber.Ctx) error {
	var req CreateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	result, err := h.invitationService.CreateInvitation(c.Context(), &req)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Created(c, result)
}

// GetInvitationDetails retrieves invitation details
// @Summary Get invitation details
// @Description Get invitation information by token
// @Tags invitations
// @Produce json
// @Param token path string true "Invitation token"
// @Success 200 {object} response.SuccessResponse{data=InvitationResponse}
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/invite/{token} [get]
func (h *Handler) GetInvitationDetails(c *fiber.Ctx) error {
	token := c.Params("token")

	result, err := h.invitationService.GetInvitationDetails(c.Context(), token)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, result)
}

// CompleteInvitation completes user registration
// @Summary Complete invitation
// @Description Complete user registration with username and password
// @Tags invitations
// @Accept json
// @Produce json
// @Param token path string true "Invitation token"
// @Param credentials body CompleteInvitationRequest true "User credentials"
// @Success 200 {object} response.SuccessResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/invite/{token}/complete [post]
func (h *Handler) CompleteInvitation(c *fiber.Ctx) error {
	token := c.Params("token")

	var req CompleteInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, err)
	}

	if err := h.invitationService.CompleteInvitation(c.Context(), token, &req); err != nil {
		return response.Error(c, err)
	}

	return response.SuccessWithMessage(c, "registration completed successfully", nil)
}
