package api

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/config"
	"github.com/skndash96/lastnight-backend/internal/handler"
	"github.com/skndash96/lastnight-backend/internal/provider"
	"github.com/skndash96/lastnight-backend/internal/repository"
	"github.com/skndash96/lastnight-backend/internal/service"
)

func RegisterRoutes(e *echo.Echo, cfg *config.AppConfig, pool *pgxpool.Pool) {
	r := e.Group("/api")

	authRepo := repository.NewAuthRepository(pool)
	teamRepo := repository.NewTeamRepository(pool)

	sessionProvider := provider.NewSessionProvider(cfg.Auth.Session, authRepo)
	uploadProvider, err := provider.NewUploadProvider(cfg.Minio)
	if err != nil {
		log.Fatalf("failed to initialize upload provider: %v", err)
	}

	r.Use(auth.SessionMW(sessionProvider, cfg.Auth.Cookie))

	{
		h := handler.NewHealthHandler()
		g := r.Group("/health")
		g.GET("", h.HealthCheck)
	}

	{
		authSrv := service.NewAuthService(pool, sessionProvider)

		h := handler.NewAuthHandler(cfg, authSrv)
		g := r.Group("/auth")
		g.POST("/login", h.Login)
		g.POST("/register", h.Register)
		g.DELETE("/logout", h.Logout)
	}

	{
		teamSrv := service.NewTeamService(pool)
		tagSrv := service.NewTagService(pool)

		team_h := handler.NewTeamHandler(teamSrv)
		tag_h := handler.NewTagHandler(tagSrv)

		teamsG := r.Group("/teams")

		teamsG.GET("", team_h.GetTeams)
		teamsG.POST("/default", team_h.JoinDefaultTeam)

		teamG := teamsG.Group("/:teamID")
		teamG.Use(auth.TeamMW(teamRepo))

		teamG.GET("/filters", tag_h.ListFilters)
		teamG.PUT("/filters", tag_h.UpdateFilters)

		teamG.POST("/tags", tag_h.CreateTagKey)
		teamG.PUT("/tags/:tagID", tag_h.UpdateTagKey)
		teamG.DELETE("/tags/:tagID", tag_h.DeleteTagKey)

		teamG.POST("/tags/:tagID/values", tag_h.CreateTagValue)
		teamG.DELETE("/tags/:tagID/values/:tagValueID", tag_h.DeleteTagValue)

		{
			uploadSrv := service.NewUploadService(uploadProvider, pool)
			h := handler.NewUploadHandler(uploadSrv)

			uploadsG := teamG.Group("/uploads")
			uploadsG.POST("/presign", h.PresignUpload)
			uploadsG.POST("/complete", h.CompleteUpload)
		}
	}
}
