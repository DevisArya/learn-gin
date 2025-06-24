package helper

import (
	"fmt"
)

// Custom error struct untuk validasi
type CustomErrors struct {
	StatusCode int      `json:"statusCode"`
	Messages   []string `json:"errors"`
}

func (e *CustomErrors) Error() string {
	return fmt.Sprintf("error: %v", e.Messages)
}

// Fungsi untuk membuat error validasi
func NewCustomError(statusCode int, messages []string) *CustomErrors {
	return &CustomErrors{
		StatusCode: statusCode,
		Messages:   messages,
	}
}
