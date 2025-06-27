package helper

import (
	"context"
	"net/http"

	"github.com/DevisArya/BE-challenge-syn/NotificationService/internal/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// send success response
func SendSuccessResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

// send success response with data
func SendSuccessResponseWithData(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"data":    data,
	})
}

// send success response with data + pagination
func SendSuccessResponseWithPagination(c *gin.Context, statusCode int, message string, data interface{}, paging *dto.PaginationResponse) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"data":    data,
		"pagination": map[string]int{
			"total_page":   paging.TotalPage,
			"total_record": paging.TotalRecord,
			"limit":        paging.Limit,
			"current_page": paging.CurrentPage,
		},
	})
}

// send error response
func SendErrorResponse(c *gin.Context, statusCode int, errors []string) {
	c.JSON(statusCode, gin.H{
		"errors": errors,
	})
}

// send custom error response
func SendCustomErrorResponse(c *gin.Context, err error) {

	// Cek error validasi
	if customErr, ok := err.(*CustomErrors); ok {
		SendErrorResponse(c, customErr.StatusCode, customErr.Messages)
	}

	// Cek error timeout
	if err == context.DeadlineExceeded {
		SendErrorResponse(c, http.StatusRequestTimeout, []string{"request timed out"})
	}

	//cek error not found
	if err == gorm.ErrRecordNotFound {
		SendErrorResponse(c, http.StatusNotFound, []string{"record not found"})
	}

	// Untuk error lainnya
	SendErrorResponse(c, http.StatusInternalServerError, []string{"internal server error"})
}
