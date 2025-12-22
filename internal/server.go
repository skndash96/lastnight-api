package api

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skndash96/lastnight-backend/internal/config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Server() error {
	_ = godotenv.Load()

	ctx := context.Background()

	appCfg := config.New()

	pool, err := pgxpool.New(ctx, appCfg.DbURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	_, err =  pool.Query(context.Background(), "SELECT 1;");
	if err != nil {
		return err
	}

	e := echo.New()

	e.Validator = NewCustomValidator()
	e.HTTPErrorHandler = ErrorHandler

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlash())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	RegisterRoutes(e, appCfg, pool)

	err = e.Start(fmt.Sprintf("localhost:%d", appCfg.Port))
	if err != nil {
		return err
	}

	return nil
}
