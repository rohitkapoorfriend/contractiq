package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/contractiq/contractiq/pkg/apperror"
)

// JSON sends a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

// OK sends a 200 response.
func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// Created sends a 201 response.
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// NoContent sends a 204 response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ErrorResponse is the standard error response format.
type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

// ErrorBody contains error details.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error sends an appropriate error response based on the error type.
func Error(w http.ResponseWriter, err error) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		JSON(w, appErr.HTTPStatus(), ErrorResponse{
			Error: ErrorBody{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			},
		})
		return
	}

	JSON(w, http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeInternal),
			Message: "an internal error occurred",
		},
	})
}

// Paginated is a standard paginated response.
type Paginated struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination metadata for list endpoints.
type Pagination struct {
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}
