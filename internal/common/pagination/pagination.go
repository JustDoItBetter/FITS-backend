// Package pagination provides utilities for handling paginated API responses.
package pagination

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	// DefaultPage is the default page number when not specified
	DefaultPage = 1
	// DefaultLimit is the default number of items per page
	DefaultLimit = 20
	// MaxLimit is the maximum allowed items per page to prevent performance issues
	MaxLimit = 100
)

// Params holds pagination parameters extracted from query string
type Params struct {
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"20"`
}

// Offset calculates the database offset from page number and limit
// Enables efficient OFFSET-LIMIT queries without loading all records
func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

// Response wraps paginated data with metadata for client navigation
type Response struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page" example:"1"`
	Limit      int         `json:"limit" example:"20"`
	TotalCount int64       `json:"total_count" example:"150"`
	TotalPages int         `json:"total_pages" example:"8"`
	Success    bool        `json:"success" example:"true"`
}

// ExtractParams parses pagination parameters from Fiber context query string
// Applies default values and validates limits to prevent abuse
func ExtractParams(c *fiber.Ctx) Params {
	page := DefaultPage
	limit := DefaultLimit

	// Parse page parameter with fallback to default
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit parameter with max limit enforcement
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			// Enforce maximum limit to prevent resource exhaustion
			if l > MaxLimit {
				limit = MaxLimit
			} else {
				limit = l
			}
		}
	}

	return Params{
		Page:  page,
		Limit: limit,
	}
}

// NewResponse creates a paginated response with calculated metadata
// totalCount should be obtained from COUNT(*) query before data fetch
func NewResponse(data interface{}, params Params, totalCount int64) Response {
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.Limit)))

	// Ensure at least 1 page even with no results
	if totalPages == 0 {
		totalPages = 1
	}

	return Response{
		Data:       data,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
		Success:    true,
	}
}
