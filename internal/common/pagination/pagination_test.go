package pagination

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestParams_Offset(t *testing.T) {
	tests := []struct {
		name   string
		params Params
		want   int
	}{
		{
			name:   "first page with default limit",
			params: Params{Page: 1, Limit: 20},
			want:   0,
		},
		{
			name:   "second page with default limit",
			params: Params{Page: 2, Limit: 20},
			want:   20,
		},
		{
			name:   "third page with custom limit",
			params: Params{Page: 3, Limit: 50},
			want:   100,
		},
		{
			name:   "page 10 with limit 10",
			params: Params{Page: 10, Limit: 10},
			want:   90,
		},
		{
			name:   "page 1 with limit 1",
			params: Params{Page: 1, Limit: 1},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.params.Offset()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExtractParams(t *testing.T) {
	tests := []struct {
		name       string
		queryPage  string
		queryLimit string
		wantPage   int
		wantLimit  int
	}{
		{
			name:       "no query parameters - use defaults",
			queryPage:  "",
			queryLimit: "",
			wantPage:   DefaultPage,
			wantLimit:  DefaultLimit,
		},
		{
			name:       "valid page and limit",
			queryPage:  "2",
			queryLimit: "50",
			wantPage:   2,
			wantLimit:  50,
		},
		{
			name:       "page only - use default limit",
			queryPage:  "5",
			queryLimit: "",
			wantPage:   5,
			wantLimit:  DefaultLimit,
		},
		{
			name:       "limit only - use default page",
			queryPage:  "",
			queryLimit: "30",
			wantPage:   DefaultPage,
			wantLimit:  30,
		},
		{
			name:       "invalid page - use default",
			queryPage:  "invalid",
			queryLimit: "25",
			wantPage:   DefaultPage,
			wantLimit:  25,
		},
		{
			name:       "invalid limit - use default",
			queryPage:  "3",
			queryLimit: "invalid",
			wantPage:   3,
			wantLimit:  DefaultLimit,
		},
		{
			name:       "negative page - use default",
			queryPage:  "-1",
			queryLimit: "20",
			wantPage:   DefaultPage,
			wantLimit:  20,
		},
		{
			name:       "zero page - use default",
			queryPage:  "0",
			queryLimit: "20",
			wantPage:   DefaultPage,
			wantLimit:  20,
		},
		{
			name:       "negative limit - use default",
			queryPage:  "2",
			queryLimit: "-10",
			wantPage:   2,
			wantLimit:  DefaultLimit,
		},
		{
			name:       "zero limit - use default",
			queryPage:  "2",
			queryLimit: "0",
			wantPage:   2,
			wantLimit:  DefaultLimit,
		},
		{
			name:       "limit exceeds maximum - cap at MaxLimit",
			queryPage:  "1",
			queryLimit: "500",
			wantPage:   1,
			wantLimit:  MaxLimit,
		},
		{
			name:       "limit at maximum - accept",
			queryPage:  "1",
			queryLimit: "100",
			wantPage:   1,
			wantLimit:  MaxLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fiber context with query parameters
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			// Build query string
			if tt.queryPage != "" {
				ctx.Request().URI().QueryArgs().Add("page", tt.queryPage)
			}
			if tt.queryLimit != "" {
				ctx.Request().URI().QueryArgs().Add("limit", tt.queryLimit)
			}

			// Extract params
			params := ExtractParams(ctx)

			// Assert
			assert.Equal(t, tt.wantPage, params.Page, "page mismatch")
			assert.Equal(t, tt.wantLimit, params.Limit, "limit mismatch")
		})
	}
}

func TestNewResponse(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		params         Params
		totalCount     int64
		wantPage       int
		wantLimit      int
		wantTotalCount int64
		wantTotalPages int
	}{
		{
			name:           "normal case - multiple pages",
			data:           []string{"item1", "item2", "item3"},
			params:         Params{Page: 1, Limit: 20},
			totalCount:     100,
			wantPage:       1,
			wantLimit:      20,
			wantTotalCount: 100,
			wantTotalPages: 5,
		},
		{
			name:           "single page - all data fits",
			data:           []string{"item1", "item2"},
			params:         Params{Page: 1, Limit: 20},
			totalCount:     2,
			wantPage:       1,
			wantLimit:      20,
			wantTotalCount: 2,
			wantTotalPages: 1,
		},
		{
			name:           "empty data - zero count",
			data:           []string{},
			params:         Params{Page: 1, Limit: 20},
			totalCount:     0,
			wantPage:       1,
			wantLimit:      20,
			wantTotalCount: 0,
			wantTotalPages: 1,
		},
		{
			name:           "second page",
			data:           []string{"item21", "item22"},
			params:         Params{Page: 2, Limit: 20},
			totalCount:     50,
			wantPage:       2,
			wantLimit:      20,
			wantTotalCount: 50,
			wantTotalPages: 3,
		},
		{
			name:           "partial last page",
			data:           []string{"item41", "item42", "item43"},
			params:         Params{Page: 3, Limit: 20},
			totalCount:     43,
			wantPage:       3,
			wantLimit:      20,
			wantTotalCount: 43,
			wantTotalPages: 3,
		},
		{
			name:           "exact division - no partial page",
			data:           []string{"item1"},
			params:         Params{Page: 1, Limit: 10},
			totalCount:     100,
			wantPage:       1,
			wantLimit:      10,
			wantTotalCount: 100,
			wantTotalPages: 10,
		},
		{
			name:           "small limit - many pages",
			data:           []string{"item1"},
			params:         Params{Page: 1, Limit: 1},
			totalCount:     100,
			wantPage:       1,
			wantLimit:      1,
			wantTotalCount: 100,
			wantTotalPages: 100,
		},
		{
			name:           "large limit - single page",
			data:           []string{"item1", "item2"},
			params:         Params{Page: 1, Limit: 100},
			totalCount:     50,
			wantPage:       1,
			wantLimit:      100,
			wantTotalCount: 50,
			wantTotalPages: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewResponse(tt.data, tt.params, tt.totalCount)

			assert.Equal(t, tt.data, response.Data, "data mismatch")
			assert.Equal(t, tt.wantPage, response.Page, "page mismatch")
			assert.Equal(t, tt.wantLimit, response.Limit, "limit mismatch")
			assert.Equal(t, tt.wantTotalCount, response.TotalCount, "total_count mismatch")
			assert.Equal(t, tt.wantTotalPages, response.TotalPages, "total_pages mismatch")
		})
	}
}

func TestNewResponse_TotalPagesCalculation(t *testing.T) {
	tests := []struct {
		name           string
		limit          int
		totalCount     int64
		wantTotalPages int
	}{
		{
			name:           "0 total count - should return 1 page",
			limit:          20,
			totalCount:     0,
			wantTotalPages: 1,
		},
		{
			name:           "1 item - should return 1 page",
			limit:          20,
			totalCount:     1,
			wantTotalPages: 1,
		},
		{
			name:           "exactly 20 items with limit 20 - should return 1 page",
			limit:          20,
			totalCount:     20,
			wantTotalPages: 1,
		},
		{
			name:           "21 items with limit 20 - should return 2 pages",
			limit:          20,
			totalCount:     21,
			wantTotalPages: 2,
		},
		{
			name:           "99 items with limit 20 - should return 5 pages",
			limit:          20,
			totalCount:     99,
			wantTotalPages: 5,
		},
		{
			name:           "100 items with limit 20 - should return 5 pages",
			limit:          20,
			totalCount:     100,
			wantTotalPages: 5,
		},
		{
			name:           "101 items with limit 20 - should return 6 pages",
			limit:          20,
			totalCount:     101,
			wantTotalPages: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewResponse([]string{}, Params{Page: 1, Limit: tt.limit}, tt.totalCount)
			assert.Equal(t, tt.wantTotalPages, response.TotalPages)
		})
	}
}
