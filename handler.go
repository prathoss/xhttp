package xhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

var HandlerLogErr = func(ctx context.Context, err error) {
	slog.ErrorContext(ctx, "response could not be written", slog.String("error", err.Error()))
}

// HttpHandler implements http.Handler
//
// The `HttpHandler` function should return a response model and an error.
// If an error is returned, it is handled according to its type:
//   - If the error implements the HttpProblemWriter interface, the HttpProblemWriter.WriteProblem method is called to write the error response.
//   - If the error is a generic error, an InternalServerError is created using the error and its InternalServerError.WriteProblem method is called.
//
// If the response model is nil, a "204 No Content" response is sent.
// Otherwise, the response model is encoded as JSON and sent as the response body.
//
// Example usage:
//
//	mux := http.NewServeMux()
//	mux.Handle("GET /nocontent", HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
//		return nil, nil
//	}))
type HttpHandler func(w http.ResponseWriter, r *http.Request) (any, error)

func (f HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseModel, err := f(w, r)
	if err != nil {
		if problemWriter, ok := err.(HttpProblemWriter); ok {
			if err := problemWriter.WriteProblem(r.Context(), w); err != nil {
				HandlerLogErr(r.Context(), fmt.Errorf("could not write response: %w", err))
			}
		} else {
			internalServerError := NewInternalServerError(err)
			err := internalServerError.WriteProblem(r.Context(), w)
			if err != nil {
				HandlerLogErr(r.Context(), fmt.Errorf("could not write response: %w", err))
			}
		}
		return
	}

	if responseModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseModel); err != nil {
		HandlerLogErr(r.Context(), fmt.Errorf("could not encode response body: %w", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type HttpProblemWriter interface {
	WriteProblem(ctx context.Context, w http.ResponseWriter) error
}
