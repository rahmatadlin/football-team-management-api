package main

import (
	"context"
	"log"
	"time"

	"football-team-management-api/config"
	"football-team-management-api/models"
	"football-team-management-api/utils/password"
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

	ctx := context.Background()

	hash, err := password.Hash(cfg.AdminPassword)
	if err != nil {
		log.Fatal(err)
	}
	var count int64
	db.Model(&models.Admin{}).Where("email = ?", cfg.AdminEmail).Count(&count)
	if count == 0 {
		a := &models.Admin{Email: cfg.AdminEmail, PasswordHash: hash}
		if err := db.WithContext(ctx).Create(a).Error; err != nil {
			log.Fatalf("seed admin: %v", err)
		}
		log.Printf("created admin %s", cfg.AdminEmail)
	} else {
		log.Printf("admin %s already exists, skipping", cfg.AdminEmail)
	}

	var teams []models.Team
	db.Find(&teams)
	if len(teams) >= 2 {
		log.Println("sample teams already present, skipping demo data")
		return
	}

	t1 := &models.Team{Name: "FC Merah", LogoURL: "https://example.com/logo1.png", FoundedYear: 2010, Address: "Jl. Sepak Bola 1", City: "Jakarta"}
	t2 := &models.Team{Name: "FC Biru", LogoURL: "https://example.com/logo2.png", FoundedYear: 2012, Address: "Jl. Stadion 2", City: "Bandung"}
	if err := db.WithContext(ctx).Create(t1).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.WithContext(ctx).Create(t2).Error; err != nil {
		log.Fatal(err)
	}

	p1 := &models.Player{TeamID: t1.ID, Name: "Andi Striker", Height: 175, Weight: 70, Position: models.PositionStriker, JerseyNumber: 9}
	p2 := &models.Player{TeamID: t1.ID, Name: "Budi Keeper", Height: 182, Weight: 78, Position: models.PositionGoalkeeper, JerseyNumber: 1}
	p3 := &models.Player{TeamID: t2.ID, Name: "Citra Striker", Height: 170, Weight: 65, Position: models.PositionStriker, JerseyNumber: 10}
	p4 := &models.Player{TeamID: t2.ID, Name: "Dedi Defender", Height: 178, Weight: 74, Position: models.PositionDefender, JerseyNumber: 4}
	for _, p := range []*models.Player{p1, p2, p3, p4} {
		if err := db.WithContext(ctx).Create(p).Error; err != nil {
			log.Fatal(err)
		}
	}

	d1, _ := time.ParseInLocation("2006-01-02", "2026-03-01", time.Local)
	d2, _ := time.ParseInLocation("2006-01-02", "2026-03-15", time.Local)
	m1 := &models.Match{MatchDate: d1, MatchTime: "15:00:00", HomeTeamID: t1.ID, AwayTeamID: t2.ID}
	m2 := &models.Match{MatchDate: d2, MatchTime: "16:30:00", HomeTeamID: t2.ID, AwayTeamID: t1.ID}
	if err := db.WithContext(ctx).Create(m1).Error; err != nil {
		log.Fatal(err)
	}
	if err := db.WithContext(ctx).Create(m2).Error; err != nil {
		log.Fatal(err)
	}

	res1 := &models.MatchResult{MatchID: m1.ID, HomeScore: 2, AwayScore: 1}
	if err := db.WithContext(ctx).Create(res1).Error; err != nil {
		log.Fatal(err)
	}
	goals := []models.Goal{
		{MatchID: m1.ID, PlayerID: p1.ID, GoalTime: 12},
		{MatchID: m1.ID, PlayerID: p1.ID, GoalTime: 67},
		{MatchID: m1.ID, PlayerID: p3.ID, GoalTime: 55},
	}
	for i := range goals {
		if err := db.WithContext(ctx).Create(&goals[i]).Error; err != nil {
			log.Fatal(err)
		}
	}

	log.Println("seed completed: teams, players, one finished match with goals, one scheduled match")
}
