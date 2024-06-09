package xhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProblemDetail struct {
	Status int    `json:"status"`
	Type   string `json:"type"`
	Title  string `json:"title"`
}

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ValidationProblemDetail struct {
	ProblemDetail
	InvalidParams []InvalidParam `json:"invalid-params"`
}

var _ error = (*BadRequestError)(nil)
var _ HttpProblemWriter = (*BadRequestError)(nil)

func NewBadRequestError(invalidParams ...InvalidParam) *BadRequestError {
	return &BadRequestError{invalidParams: invalidParams}
}

type BadRequestError struct {
	invalidParams []InvalidParam
}

func (b *BadRequestError) Error() string {
	return fmt.Sprintf("%v", b.invalidParams)
}

func (b *BadRequestError) WriteProblem(_ context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ValidationProblemDetail{
		ProblemDetail: ProblemDetail{
			Status: http.StatusBadRequest,
			Type:   "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.1",
			Title:  "Request parameters did not validate",
		},
		InvalidParams: b.invalidParams,
	}
	return json.NewEncoder(w).Encode(detail)
}

var _ error = (*UnauthorizedError)(nil)
var _ HttpProblemWriter = (*UnauthorizedError)(nil)

type UnauthorizedError struct {
	message string
}

func NewUnauthorizedError(message string) *UnauthorizedError {
	return &UnauthorizedError{message: message}
}

func (f *UnauthorizedError) WriteProblem(ctx context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ProblemDetail{
		Status: http.StatusUnauthorized,
		Type:   "https://datatracker.ietf.org/doc/html/rfc7235#section-3.1",
		Title:  f.message,
	}
	return json.NewEncoder(w).Encode(detail)
}

func (f *UnauthorizedError) Error() string {
	return f.message
}

var _ error = (*ForbiddenError)(nil)
var _ HttpProblemWriter = (*ForbiddenError)(nil)

type ForbiddenError struct {
	message string
}

func NewForbiddenError(message string) *ForbiddenError {
	return &ForbiddenError{message: message}
}

func (f *ForbiddenError) WriteProblem(ctx context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ProblemDetail{
		Status: http.StatusForbidden,
		Type:   "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.3",
		Title:  f.message,
	}
	return json.NewEncoder(w).Encode(detail)
}

func (f *ForbiddenError) Error() string {
	return f.message
}

var _ error = &NotFoundError{}
var _ HttpProblemWriter = &NotFoundError{}

func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		message: message,
	}
}

type NotFoundError struct {
	message string
}

func (n *NotFoundError) Error() string {
	return n.message
}

func (n *NotFoundError) WriteProblem(ctx context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ProblemDetail{
		Status: http.StatusNotFound,
		Type:   "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.4",
		Title:  n.message,
	}
	return json.NewEncoder(w).Encode(detail)
}

var _ error = (*UnsupportedMediaType)(nil)
var _ HttpProblemWriter = (*UnsupportedMediaType)(nil)

type UnsupportedMediaType struct {
}

func NewUnsupportedMediaType() *UnsupportedMediaType {
	return &UnsupportedMediaType{}
}

func (u *UnsupportedMediaType) WriteProblem(ctx context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusUnsupportedMediaType)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ProblemDetail{
		Status: http.StatusForbidden,
		Type:   "https://datatracker.ietf.org/doc/html/rfc7231#section-6.5.13",
		Title:  "unsupported media type",
	}
	return json.NewEncoder(w).Encode(detail)
}

func (u *UnsupportedMediaType) Error() string {
	return "unsupported media type"
}

var _ error = (*InternalServerError)(nil)
var _ HttpProblemWriter = (*InternalServerError)(nil)

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{
		innerError: err,
	}
}

type InternalServerError struct {
	innerError error
}

func (i *InternalServerError) Error() string {
	return i.innerError.Error()
}

func (i *InternalServerError) WriteProblem(_ context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ProblemDetail{
		Status: http.StatusInternalServerError,
		Type:   "https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.1",
		Title:  "Internal Server Error",
	}
	return json.NewEncoder(w).Encode(detail)
}

var _ error = (*ServiceUnavailableError)(nil)
var _ HttpProblemWriter = (*ServiceUnavailableError)(nil)

func NewServiceUnavailableError(err error) *ServiceUnavailableError {
	return &ServiceUnavailableError{
		innerError: err,
	}
}

type ServiceUnavailableError struct {
	innerError error
}

func (s *ServiceUnavailableError) WriteProblem(ctx context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Header().Set("Content-Type", "application/problem+json")
	detail := ProblemDetail{
		Status: http.StatusServiceUnavailable,
		Type:   "https://datatracker.ietf.org/doc/html/rfc7231#section-6.6.4",
		Title:  "The server is unavailable",
	}
	return json.NewEncoder(w).Encode(detail)
}

func (s *ServiceUnavailableError) Error() string {
	return s.innerError.Error()
}
