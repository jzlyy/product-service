package utils

import (
	"product-service/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePagination(c *gin.Context) models.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	return models.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func CalculateTotalPages(total int, pageSize int) int {
	if total == 0 {
		return 1
	}
	pages := total / pageSize
	if total%pageSize > 0 {
		pages++
	}
	return pages
}
