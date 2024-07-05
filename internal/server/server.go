package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"

	"github.com/wopoczynski/todoapp/internal/database"
	handler "github.com/wopoczynski/todoapp/internal/handlers"
	"github.com/wopoczynski/todoapp/internal/initialize"

	_ "github.com/wopoczynski/todoapp/docs"
)

type Config struct {
	HTTPServerPort       string `env:"HTTP_SERVER_PORT,default=8123"`
	*initialize.DBConfig `env:", prefix=DB_"`
}

func Run(ctx context.Context) error {
	var cfg Config

	if err := envconfig.Process(ctx, &cfg); err != nil {
		return fmt.Errorf("parsing configuration: %w", err)
	}

	db, err := initialize.DB(*cfg.DBConfig)
	defer func() {
		conn, dbErr := db.DB()
		if dbErr != nil {
			log.Err(err).Msg("db closing")
		}

		conn.Close()
	}()
	if err != nil {
		return fmt.Errorf("unable to connect: %w", err)
	}
	err = initialize.Automigrate(db)
	if err != nil {
		return fmt.Errorf("unable to migrate schema: %w", err)
	}
	r := database.NewMysqlTodoRepository(db)
	h := handler.NewHandler(r)

	e := echo.New()
	e.HideBanner = true
	todo := e.Group("/todos")
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORS())

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(_ echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")
			return nil
		},
		Skipper: func(c echo.Context) bool {
			return strings.EqualFold(c.Request().URL.Path, "/") ||
				strings.EqualFold(c.Request().URL.Path, "/ping")
		},
	}))
	e.GET("/", handler.Ping)
	e.GET("/ping", handler.Ping)

	todo.GET("", h.GetAllTodos)
	todo.POST("", h.CreateTodo)
	todo.DELETE("", h.DeleteAllTodos)
	todo.GET("/:id", h.GetTodo)
	todo.DELETE("/:id", h.DeleteTodo)
	todo.PATCH("/:id", h.UpdateTodo)

	go func() {
		err = e.Start(":" + cfg.HTTPServerPort)
		if errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("Server shutdown")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	const shutdownTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	err = e.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("shutting down server...")
	}

	log.Info().Msg("server stopped gracefully")
	return nil
}
