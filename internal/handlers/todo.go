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

// GetAllTodos godoc
//
//	@Summary        GET all todos
//	@Description    Get list of all todos created
//	@Tags           Todo App
//	@Produce        json
//	@Success        200 {object}    []database.TodoModel
//	@Router         /todos [get]
func (h *Handler) GetAllTodos(c echo.Context) error {
	ctx := c.Request().Context()
	return c.JSON(http.StatusOK, h.repository.GetAll(ctx))
}

// CreateTodo godoc
//
//	@Summary        POST todo
//	@Description    Create todo
//	@Tags           Todo App
//	@Produce        json
//	@Param          request	body	database.TodoModel	true	"Request body"
//	@Success        201
//	@Router         /todos [post]
func (h *Handler) CreateTodo(c echo.Context) error {
	ctx := c.Request().Context()
	var todo database.TodoModel

	rsp, err := h.repository.Create(ctx, &todo)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid data provided")
	}
	return c.JSON(http.StatusCreated, rsp)
}

// DeleteAllTodos godoc
//
//	@Summary        Delete todos
//	@Description    Delete all created todo by id
//	@Tags           Todo App
//	@Success        204
//	@Failure        400
//	@Router         /todos [delete]
func (h *Handler) DeleteAllTodos(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.repository.DeleteAll(ctx)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "unable to purge todos")
	}
	return c.JSON(http.StatusOK, nil)
}

// GetTodo godoc
//
//	@Summary        Get todo
//	@Description    Get created todo by id
//	@Tags           Todo App
//	@Param          id  path    string  true    "TODO ID"
//	@Success        201
//	@Failure        400
//	@Router         /todos [get]
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

// DeleteTodo godoc
//
//	@Summary        Delete todo
//	@Description    Delete created todo by id
//	@Tags           Todo App
//	@Param          id  path    string  true    "TODO ID"
//	@Success        204
//	@Failure        400
//	@Router         /todos [delete]
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

// UpdateTodo godoc
//
//	@Summary        PATCH todo
//	@Description    Update todo
//	@Tags           Todo App
//	@Produce        json
//	@Param          request	body	database.TodoModel	true	"Request body"
//	@Success        200
//	@Router         /todos [patch]
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
