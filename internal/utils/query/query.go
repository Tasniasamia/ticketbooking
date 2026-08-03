// internal/utils/query/query.go
package query;

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Params struct {
	Page    int
	Limit   int
	Search  string
	SortBy  string
	SortDir string // asc / desc

	// Dynamic filters (key = column, value = value)
	Filters map[string]string

	// Range filters
	MinPrice *int
	MaxPrice *int
	FromDate *string
	ToDate   *string
}

func Parse(c *echo.Context) Params {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	sortBy := c.QueryParam("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortDir := strings.ToLower(c.QueryParam("sort_dir"))
	if sortDir != "asc" {
		sortDir = "desc"
	}

	filters := map[string]string{}
	for key, values := range c.QueryParams() {
		if len(values) == 0 {
			continue
		}
		// skip known params
		switch key {
		case "page", "limit", "search", "sort_by", "sort_dir",
			"min_price", "max_price", "from_date", "to_date":
			continue
		}
		filters[key] = values[0]
	}

	var minPrice, maxPrice *int
	if v := c.QueryParam("min_price"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			minPrice = &n
		}
	}
	if v := c.QueryParam("max_price"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			maxPrice = &n
		}
	}

	var fromDate, toDate *string
	if v := c.QueryParam("from_date"); v != "" {
		fromDate = &v
	}
	if v := c.QueryParam("to_date"); v != "" {
		toDate = &v
	}

	return Params{
		Page:     page,
		Limit:    limit,
		Search:   c.QueryParam("search"),
		SortBy:   sortBy,
		SortDir:  sortDir,
		Filters:  filters,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		FromDate: fromDate,
		ToDate:   toDate,
	}
}

// Apply → GORM query-তে সব কিছু apply করে
func Apply(db *gorm.DB, p Params, searchFields []string) *gorm.DB {
	// 1. Search (ILIKE for postgres)
	if p.Search != "" && len(searchFields) > 0 {
		var conditions []string
		var args []interface{}
		for _, field := range searchFields {
			conditions = append(conditions, field+" ILIKE ?")
			args = append(args, "%"+p.Search+"%")
		}
		db = db.Where(strings.Join(conditions, " OR "), args...)
	}

	// 2. Exact filters (status=active, category=vip ইত্যাদি)
	for col, val := range p.Filters {
		db = db.Where(col+" = ?", val)
	}

	// 3. Price range
	if p.MinPrice != nil {
		db = db.Where("price >= ?", *p.MinPrice)
	}
	if p.MaxPrice != nil {
		db = db.Where("price <= ?", *p.MaxPrice)
	}

	// 4. Date range
	if p.FromDate != nil {
		db = db.Where("starts_at >= ?", *p.FromDate)
	}
	if p.ToDate != nil {
		db = db.Where("starts_at <= ?", *p.ToDate)
	}

	// 5. Sorting
	db = db.Order(p.SortBy + " " + p.SortDir)

	return db
}

func Paginate(db *gorm.DB, p Params) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return db.Limit(p.Limit).Offset(offset)
}