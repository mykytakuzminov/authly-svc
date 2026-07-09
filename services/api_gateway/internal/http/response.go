package http

import (
	"net/http"
	"strconv"
)

type SuccessResponse struct {
	Data any `json:"data"`
}
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func responseBadRequest(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusBadRequest, ErrorResponse{
		Error: &APIError{
			Code:    http.StatusText(http.StatusBadRequest),
			Message: msg,
		},
	})
}

func responseUnauthorized(w http.ResponseWriter, msg string) {
	WriteJSON(w, http.StatusUnauthorized, ErrorResponse{
		Error: &APIError{
			Code:    http.StatusText(http.StatusUnauthorized),
			Message: msg,
		},
	})
}

func responseError(w http.ResponseWriter, code int, msg string) {
	WriteJSON(w, code, ErrorResponse{
		Error: &APIError{
			Code:    strconv.Itoa(code),
			Message: msg,
		},
	})
}

func responseNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func responseSuccess(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, SuccessResponse{
		Data: data,
	})
}
