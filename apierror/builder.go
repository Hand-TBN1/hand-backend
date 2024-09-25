package apierror

import (
	"net/http"
)

// ApiError represents a standard API error response.
type ApiError struct {
	HttpStatus int         `json:"status"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload,omitempty"`
}

// ApiErrorBuilder is used to build ApiError step by step.
type ApiErrorBuilder struct {
	httpStatus int
	message    string
	payload    interface{}
}

// NewApiErrorBuilder initializes a new ApiErrorBuilder with default values.
func NewApiErrorBuilder() *ApiErrorBuilder {
	return &ApiErrorBuilder{
		httpStatus: http.StatusInternalServerError, // default status
	}
}

// WithStatus sets the HTTP status for the error.
func (b *ApiErrorBuilder) WithStatus(status int) *ApiErrorBuilder {
	b.httpStatus = status
	return b
}

// WithMessage sets the error message.
func (b *ApiErrorBuilder) WithMessage(message string) *ApiErrorBuilder {
	b.message = message
	return b
}

// WithPayload adds additional payload to the error (optional).
func (b *ApiErrorBuilder) WithPayload(payload interface{}) *ApiErrorBuilder {
	b.payload = payload
	return b
}

// Build creates the ApiError instance.
func (b *ApiErrorBuilder) Build() *ApiError {
	return &ApiError{
		HttpStatus: b.httpStatus,
		Message:    b.message,
		Payload:    b.payload,
	}
}
