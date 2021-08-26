package api

import (
	"net/http"
)

// APIResponse is format return
type APIResponse struct {
	StatusCode int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

var (
	RequestOK  = NewResponse(http.StatusOK, http.StatusText(http.StatusOK), "")
	EmptyParam = NewResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "empty param")

	BadRequest  = NewResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "")
	ServerError = NewResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "")
	NotFound    = NewResponse(http.StatusNotFound, http.StatusText(http.StatusInternalServerError), "")
)

func NewResponse(statusCode int, Msg string, Data interface{}) *APIResponse {
	return &APIResponse{
		StatusCode: statusCode,
		Msg:        Msg,
		Data:       Data,
	}
}

func NewServerError(msg string) *APIResponse {
	return NewResponse(http.StatusInternalServerError, msg, "")
}
