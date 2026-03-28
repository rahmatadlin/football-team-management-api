package routes

import (
	"github.com/gin-gonic/gin"

	"football-team-management-api/config"
	"football-team-management-api/handler"
	"football-team-management-api/middleware"
)

func Register(r *gin.Engine, cfg *config.Config, h *Handlers) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	v1.POST("/auth/login", h.Auth.Login)

	protected := v1.Group("")
	protected.Use(middleware.AdminAuth(cfg))

	protected.GET("/teams", h.Team.List)
	protected.POST("/teams", h.Team.Create)
	// Rute lebih spesifik dulu: Gin mewajibkan nama wildcard sama (:id) di seluruh /teams/:id/...
	protected.GET("/teams/:id/players", h.Player.ListByTeam)
	protected.POST("/teams/:id/players", h.Player.Create)
	protected.GET("/teams/:id", h.Team.Get)
	protected.PUT("/teams/:id", h.Team.Update)
	protected.DELETE("/teams/:id", h.Team.Delete)
	protected.PUT("/players/:id", h.Player.Update)
	protected.DELETE("/players/:id", h.Player.Delete)

	protected.GET("/matches/schedules", h.Match.ListSchedules)
	protected.POST("/matches/schedules", h.Match.CreateSchedule)
	protected.POST("/matches/:id/results", h.Match.SubmitResult)

	protected.GET("/matches/:id/report", h.Report.MatchReport)
}

type Handlers struct {
	Auth   *handler.AuthHandler
	Team   *handler.TeamHandler
	Player *handler.PlayerHandler
	Match  *handler.MatchHandler
	Report *handler.ReportHandler
}
