package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
	"football-team-management-api/repository"
	"football-team-management-api/utils/apperror"
)

type GoalInput struct {
	PlayerID string `json:"player_id" binding:"required,uuid"`
	GoalTime int    `json:"goal_time" binding:"required,min=0,max=130"`
}

type MatchService interface {
	CreateSchedule(ctx context.Context, m *models.Match) error
	ListSchedules(ctx context.Context) ([]models.Match, error)
	SubmitResult(ctx context.Context, matchID uuid.UUID, homeScore, awayScore int, goals []GoalInput) error
}

type matchService struct {
	matches  repository.MatchRepository
	teams    repository.TeamRepository
	results  repository.MatchResultRepository
	goals    repository.GoalRepository
	players  repository.PlayerRepository
}

func NewMatchService(
	matches repository.MatchRepository,
	teams repository.TeamRepository,
	results repository.MatchResultRepository,
	goals repository.GoalRepository,
	players repository.PlayerRepository,
) MatchService {
	return &matchService{
		matches: matches,
		teams:   teams,
		results: results,
		goals:   goals,
		players: players,
	}
}

func (s *matchService) CreateSchedule(ctx context.Context, m *models.Match) error {
	if m.HomeTeamID == m.AwayTeamID {
		return apperror.BadRequest("home_team_id and away_team_id must differ")
	}
	if _, err := s.teams.FindByID(ctx, m.HomeTeamID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("home team not found")
		}
		return err
	}
	if _, err := s.teams.FindByID(ctx, m.AwayTeamID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("away team not found")
		}
		return err
	}
	return s.matches.Create(ctx, m)
}

func (s *matchService) ListSchedules(ctx context.Context) ([]models.Match, error) {
	return s.matches.List(ctx)
}

func (s *matchService) SubmitResult(ctx context.Context, matchID uuid.UUID, homeScore, awayScore int, inputs []GoalInput) error {
	if homeScore < 0 || awayScore < 0 {
		return apperror.BadRequest("scores must be non-negative")
	}
	wantGoals := homeScore + awayScore
	if len(inputs) != wantGoals {
		return apperror.BadRequest("goals count must equal home_score + away_score")
	}
	match, err := s.matches.FindByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("match not found")
		}
		return err
	}
	playerIDs := make([]uuid.UUID, 0, len(inputs))
	for _, g := range inputs {
		pid, err := uuid.Parse(g.PlayerID)
		if err != nil {
			return apperror.BadRequest("invalid player_id in goals")
		}
		playerIDs = append(playerIDs, pid)
	}
	players, err := s.players.FindByIDs(ctx, playerIDs)
	if err != nil {
		return err
	}
	byID := make(map[uuid.UUID]models.Player)
	for _, p := range players {
		byID[p.ID] = p
	}
	for i, g := range inputs {
		pid := playerIDs[i]
		pl, ok := byID[pid]
		if !ok {
			return apperror.BadRequest("unknown player_id in goals: " + g.PlayerID)
		}
		if pl.TeamID != match.HomeTeamID && pl.TeamID != match.AwayTeamID {
			return apperror.BadRequest("goal scorer must belong to home or away team")
		}
	}
	goalModels := make([]models.Goal, 0, len(inputs))
	for i, g := range inputs {
		goalModels = append(goalModels, models.Goal{
			PlayerID: playerIDs[i],
			GoalTime: g.GoalTime,
		})
	}
	if err := s.results.Upsert(ctx, &models.MatchResult{
		MatchID:   matchID,
		HomeScore: homeScore,
		AwayScore: awayScore,
	}); err != nil {
		return err
	}
	return s.goals.ReplaceForMatch(ctx, matchID, goalModels)
}
