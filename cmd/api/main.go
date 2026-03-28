package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"football-team-management-api/config"
	"football-team-management-api/handler"
	"football-team-management-api/models"
	"football-team-management-api/repository"
	"football-team-management-api/routes"
	"football-team-management-api/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := config.NewDB(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Admin{},
		&models.Team{},
		&models.Player{},
		&models.Match{},
		&models.MatchResult{},
		&models.Goal{},
	); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	adminRepo := repository.NewAdminRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	resultRepo := repository.NewMatchResultRepository(db)
	goalRepo := repository.NewGoalRepository(db)

	authSvc := service.NewAuthService(cfg, adminRepo)
	teamSvc := service.NewTeamService(teamRepo)
	playerSvc := service.NewPlayerService(playerRepo, teamRepo)
	matchSvc := service.NewMatchService(matchRepo, teamRepo, resultRepo, goalRepo, playerRepo)
	reportSvc := service.NewReportService(matchRepo, goalRepo)

	h := &routes.Handlers{
		Auth:   handler.NewAuthHandler(authSvc),
		Team:   handler.NewTeamHandler(teamSvc),
		Player: handler.NewPlayerHandler(playerSvc),
		Match:  handler.NewMatchHandler(matchSvc),
		Report: handler.NewReportHandler(reportSvc),
	}

	if cfg.AppEnv != "development" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	routes.Register(r, cfg, h)

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
