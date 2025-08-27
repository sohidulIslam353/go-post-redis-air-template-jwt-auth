package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (limit, offset int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "2"))
	if err != nil || pageSize < 1 {
		pageSize = 2
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

func GetPaginationParamsWithOffset(c *gin.Context) (limit, offset int) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return
}

// Cursor Pagination Helper
func GetCursorPagination(c *gin.Context) (lastID int64, limit int) {
	// last_id → cursor
	lastIDStr := c.DefaultQuery("last_id", "0")
	lastID, _ = strconv.ParseInt(lastIDStr, 10, 64) // ✅ base 10

	// page_size → limit
	pageSizeStr := c.DefaultQuery("page_size", "10") // 10 per page
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 2
	}

	return lastID, pageSize
}

func GetCursorPaginationWithOffset(c *gin.Context) (lastID int64, limit int) {
	// last_id → cursor
	lastIDStr := c.DefaultQuery("last_id", "0")
	lastID, _ = strconv.ParseInt(lastIDStr, 10, 64) // ✅ base 10

	// page_size → limit
	pageSizeStr := c.DefaultQuery("page_size", "50") // 50 per page
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 50
	}

	return lastID, pageSize
}
