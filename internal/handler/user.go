package handler

import (
	"errors"
	"net/http"
	"time"

	"refactoring/internal/model"
	"refactoring/internal/response"
	"refactoring/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserHandler struct {
	store *store.Store
}

func NewUserHandler(store *store.Store) UserHandler {
	return UserHandler{
		store: store,
	}
}

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (c *CreateUserRequest) Bind(r *http.Request) error {
	return nil
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	request := CreateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	u := model.User{
		CreatedAt:   time.Now(),
		DisplayName: request.DisplayName,
		Email:       request.Email,
	}

	id, err := h.store.CreateUser(u)
	if err != nil {
		_ = render.Render(w, r, response.ErrInternal(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"user_id": id,
	})
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.store.GetUser(id)
	if errors.Is(err, store.ErrUserNotFound) {
		_ = render.Render(w, r, response.ErrNotFound(err))
		return
	}
	if err != nil {
		_ = render.Render(w, r, response.ErrInternal(err))
		return
	}

	render.JSON(w, r, user)
}

type UpdateUserRequest struct {
	DisplayName string `json:"display_name"`
}

func (c *UpdateUserRequest) Bind(r *http.Request) error {
	return nil
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	request := UpdateUserRequest{}

	if err := render.Bind(r, &request); err != nil {
		_ = render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}

	id := chi.URLParam(r, "id")

	user, err := h.store.GetUser(id)
	if errors.Is(err, store.ErrUserNotFound) {
		_ = render.Render(w, r, response.ErrNotFound(err))
		return
	}
	if err != nil {
		_ = render.Render(w, r, response.ErrInternal(err))
		return
	}

	user.DisplayName = request.DisplayName
	err = h.store.UpdateUser(id, user)
	if errors.Is(err, store.ErrUserNotFound) {
		_ = render.Render(w, r, response.ErrNotFound(err))
		return
	}
	if err != nil {
		_ = render.Render(w, r, response.ErrInternal(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.store.DeleteUser(id)
	if errors.Is(err, store.ErrUserNotFound) {
		_ = render.Render(w, r, response.ErrNotFound(err))
		return
	}
	if err != nil {
		_ = render.Render(w, r, response.ErrInternal(err))
		return
	}

	render.Status(r, http.StatusNoContent)
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	list := h.store.ListUsers()
	render.JSON(w, r, list)
}
