package handler

import (
	"errors"
	"net/http"
	"refactoring/internal/response"
	"refactoring/internal/store"

	"github.com/go-chi/render"
)

// ErrorHTTPHandlerFunc is equivalent to http.HandlerFunc, except it allows
// handlers to return Go errors
type ErrorHTTPHandlerFunc func(http.ResponseWriter, *http.Request) error

// ReportError is a middleware function that transforms Go errors returned from
// a `ErrorHTTPHandlerFunc` into HTTP error responses
func ReportError(handler ErrorHTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)

		if err == nil {
			return
		}

		if errors.Is(err, store.ErrUserNotFound) {
			_ = render.Render(w, r, response.ErrNotFound(err))
			return
		}

		var res *response.ErrResponse
		if errors.As(err, &res) {
			_ = render.Render(w, r, res)
			return
		}
		_ = render.Render(w, r, response.ErrInternal(response.ErrNotFound(err)))
	}
}
