package main

import (
	"net/http"
	"time"

	"refactoring/internal/handler"
	"refactoring/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const storePath = `users.json`

func main() {
	userStore, _ := store.Open(storePath)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	})

	userHandler := handler.NewUserHandler(userStore)

	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.Search)
				r.Post("/", handler.ReportError(userHandler.Create))

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", handler.ReportError(userHandler.Get))
					r.Patch("/", handler.ReportError(userHandler.Update))
					r.Delete("/", handler.ReportError(userHandler.Delete))
				})
			})
		})
	})

	http.ListenAndServe(":3333", r)
}
