package main

import (
	"encoding/json"
	"net/http"
)

// Handler is a function that can output a JSONResponse
type Handler func(*http.Request) (*JSONResponse, error)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := h(r)
	if err != nil {
		response = &JSONResponse{Code: 500, Data: err}
	}

	w.WriteHeader(response.Code)
	json.NewEncoder(w).Encode(response.Data)
}

// Route encodes a handler with a
// Method and path prefix.
type Route struct {
	Method, Prefix string
	Handler        Handler
}

type JSONResponse struct {
	Code int
	Data any
}

func OKResponse(data any) *JSONResponse {
	return &JSONResponse{Code: 200, Data: data}
}
