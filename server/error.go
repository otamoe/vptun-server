package server

import (
	"fmt"
	"net/http"
)

type (
	ValidateError struct {
		Name  string
		Value interface{}
	}
	NotFoundError struct {
		Name  string
		Value interface{}
	}
	HttpError struct {
		Err    error
		Status int
	}
)

func (validateError *ValidateError) Error() string {
	return fmt.Sprintf("validate: %s field error", validateError.Name)
}
func (notFoundError *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", notFoundError.Name)
}

func (httpError *HttpError) Error() string {
	if httpError.Err != nil {
		return httpError.Err.Error()
	}
	return httpError.StatusText()
}

func (httpError *HttpError) StatusText() string {
	code := httpError.StatusCode()
	return http.StatusText(code)
}

func (httpError *HttpError) StatusCode() int {
	if httpError.Status == 0 {
		return 500
	}
	return httpError.Status
}
