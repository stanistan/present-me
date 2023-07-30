package http

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Request = http.Request

// Handler is a function that can output a JSONResponse
type Handler func(*http.Request) (*JSONResponse, error)

// ServeHTTP fulfills the http.Handler interface for Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := h(r)
	if err != nil {
		log.Ctx(r.Context()).Error().Msgf("handler error: %s", err)
		response = ErrResponse(err)
	}

	if err := response.Write(w); err != nil {
		_ = ErrResponse(errors.New("Internal Error")).Write(w)
	}
}

// Route encodes a handler with a Method and path prefix.
type Route struct {
	Method, Prefix string
	Handler        Handler
}

func GET(prefix string, handler Handler) Route {
	return Route{
		Method:  "GET",
		Prefix:  prefix,
		Handler: handler,
	}
}

func Routes(rs ...Route) []Route {
	return rs
}

// JSONResponse represents our JSON with response code.
type JSONResponse struct {
	Code int
	Data any
}

// OKResponse
func OKResponse(data any) *JSONResponse {
	return &JSONResponse{Code: 200, Data: data}
}

func ErrResponse(err error) *JSONResponse {
	return &JSONResponse{
		Code: 500,
		Data: map[string]string{
			"msg": err.Error(),
		},
	}
}

func (r *JSONResponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	err := json.NewEncoder(w).Encode(r.Data)
	if err != nil {
		log.Error().Err(err).Msg("could not write JSON response")
		return errors.Wrap(err, "failed writing JSON")
	}

	return nil
}
