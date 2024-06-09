package xhttp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prathoss/xhttp"
)

func TestHttpHandler_ServeHTTP(t *testing.T) {
	xhttp.HandlerLogErr = func(_ context.Context, _ error) {
		// noop
	}

	tests := []struct {
		name       string
		httpFunc   xhttp.HttpHandler
		wantStatus int
	}{
		{
			name: "OK response status",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return "ok", nil
			}),
			wantStatus: http.StatusOK,
		},
		{
			name: "NoContent response when body is nil",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return nil, nil
			}),
			wantStatus: http.StatusNoContent,
		},
		{
			name: "ServiceUnavailableError",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return nil, xhttp.NewServiceUnavailableError(errors.New("service unavailable"))
			}),
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name: "BadRequestError",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return nil, xhttp.NewBadRequestError(xhttp.InvalidParam{Name: "param1", Reason: "error1"})
			}),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "InternalServerError",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return nil, xhttp.NewInternalServerError(errors.New("internal error"))
			}),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "NotFoundError",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return nil, xhttp.NewNotFoundError("not found")
			}),
			wantStatus: http.StatusNotFound,
		},
		{
			name: "generic error",
			httpFunc: xhttp.HttpHandler(func(w http.ResponseWriter, r *http.Request) (any, error) {
				return nil, errors.New("generic error")
			}),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()

			tt.httpFunc.ServeHTTP(recorder, request)

			// Assert
			if status := recorder.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatus)
			}
		})
	}
}
