package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/wopoczynski/todoapp/internal/database"
)

type Handler struct {
	repository database.TodoRepository
}

func NewHandler(r database.TodoRepository) *Handler {
	return &Handler{
		repository: r,
	}
}

func (h *Handler) GetAllTodos(c echo.Context) error {
	ctx := c.Request().Context()
	return c.JSON(http.StatusOK, h.repository.GetAll(ctx))
}

func (h *Handler) CreateTodo(c echo.Context) error {
	ctx := c.Request().Context()
	var todo database.TodoModel

	rsp, err := h.repository.Create(ctx, &todo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid data provided")
	}
	return c.JSON(http.StatusCreated, rsp)
}

func (h *Handler) DeleteAllTodos(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.repository.DeleteAll(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "unable to purge todos")
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *Handler) GetTodo(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "Missing id")
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id")
	}

	rsp, err := h.repository.Get(ctx, uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid data provided")
	}
	return c.JSON(http.StatusCreated, rsp)
}

func (h *Handler) DeleteTodo(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "Missing id")
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid id")
	}

	err = h.repository.Delete(ctx, uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid data provided")
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) UpdateTodo(c echo.Context) error {
	ctx := c.Request().Context()
	var todo database.TodoModel

	if err := c.Bind(&todo); err != nil {
		log.Error().Err(err).Msg("invalid payload for todo provided")
		return c.JSON(http.StatusBadRequest, "invalid data provided")
	}

	err := h.repository.Update(ctx, &todo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid data provided")
	}
	return c.JSON(http.StatusOK, nil)
}
