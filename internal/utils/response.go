package utils

import "time"

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Message   string      `json:"message"`
	Status    int         `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	Metadata  *Metadata   `json:"metadata,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Metadata holds additional information about the response
type Metadata struct {
	// Pagination info
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalCount int `json:"total_count,omitempty"`

	// Additional metadata fields
	Version    string                 `json:"version,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
	CustomData map[string]interface{} `json:"custom_data,omitempty"`
}

// ErrorInfo provides detailed error information
type ErrorInfo struct {
	Code       string   `json:"code,omitempty"`
	Details    string   `json:"details,omitempty"`
	Field      string   `json:"field,omitempty"`
	Validation []string `json:"validation,omitempty"`
}

// Success creates a successful response with optional metadata
func Success(message string, data interface{}, metadata ...*Metadata) StandardResponse {
	response := StandardResponse{
		Message:   message,
		Status:    200,
		Data:      data,
		Timestamp: time.Now(),
	}

	if len(metadata) > 0 {
		response.Metadata = metadata[0]
	}

	return response
}

// Error creates an error response with optional error details
func Error(message string, status int, errorInfo ...*ErrorInfo) StandardResponse {
	response := StandardResponse{
		Message:   message,
		Status:    status,
		Timestamp: time.Now(),
	}

	if len(errorInfo) > 0 {
		response.Error = errorInfo[0]
	}

	return response
}

// NewPaginationMetadata creates metadata for paginated responses
func NewPaginationMetadata(page, perPage, totalCount int) *Metadata {
	totalPages := (totalCount + perPage - 1) / perPage

	return &Metadata{
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		TotalCount: totalCount,
	}
}

// NewErrorInfo creates a new ErrorInfo instance
func NewErrorInfo(code string, details string, field string, validation []string) *ErrorInfo {
	return &ErrorInfo{
		Code:       code,
		Details:    details,
		Field:      field,
		Validation: validation,
	}
}
