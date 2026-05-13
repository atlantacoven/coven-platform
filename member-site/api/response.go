// Package api manages the low-level operations for an http api
package api

import (
	"net/http"
)

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	// TODO: number of pages, number of items
}

type apiResponse struct {
	Status     string         `json:"status"`
	Error      *errorResponse `json:"error,omitempty"`
	Data       any            `json:"data,omitempty"`
	Pagination *Pagination    `json:"pagination,omitempty"`
}

type errorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Respond[T any](w http.ResponseWriter, data *T, status string) {
	if status == "" {
		status = "OK"
	}
	r := apiResponse{
		Status: status,
		Data:   data,
	}
	MarshalBody(w, &r)
}

func RespondPage[T any](w http.ResponseWriter, data []T, p *Pagination, status string) {
	r := apiResponse{
		Status:     status,
		Data:       data,
		Pagination: p,
	}
	MarshalBody(w, &r)
}

func RespondError(w http.ResponseWriter, err error) {
	r := errorResponse{
		Message: err.Error(),
		Code:    http.StatusInternalServerError,
	}
	if herr, ok := err.(HTTPResponseError); ok {
		r.Message = herr.Error()
		r.Code = herr.Status()
	}
	if IsEnv(Production) {
		// hide the full error
		r.Message = "server error"
	}
	ar := apiResponse{
		Status: "ERROR",
		Error:  &r,
	}
	w.WriteHeader(r.Code) // send the error code
	MarshalBody(w, &ar)
}

func RespondBadFormat(w http.ResponseWriter, cause error) {
	RespondError(w, httpError{cause: cause, status: http.StatusBadRequest})
}

func UnauthorizedError(w http.ResponseWriter, cause error) {
	RespondError(w, httpError{cause: cause, status: http.StatusUnauthorized})
}

func ForbiddenError(w http.ResponseWriter, cause error) {
	RespondError(w, httpError{cause: cause, status: http.StatusForbidden})
}

func NotFoundError(w http.ResponseWriter, cause error) {
	RespondError(w, httpError{cause: cause, status: http.StatusNotFound})
}
