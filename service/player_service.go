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

type PlayerService interface {
	Create(ctx context.Context, p *models.Player) error
	Update(ctx context.Context, id uuid.UUID, patch *models.Player) (*models.Player, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]models.Player, error)
}

type playerService struct {
	players repository.PlayerRepository
	teams   repository.TeamRepository
}

func NewPlayerService(players repository.PlayerRepository, teams repository.TeamRepository) PlayerService {
	return &playerService{players: players, teams: teams}
}

func (s *playerService) ensureTeam(ctx context.Context, teamID uuid.UUID) error {
	_, err := s.teams.FindByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("team not found")
		}
		return err
	}
	return nil
}

func (s *playerService) jerseyAvailable(ctx context.Context, teamID uuid.UUID, jersey int, exclude *uuid.UUID) error {
	n, err := s.players.CountJerseyInTeam(ctx, teamID, jersey, exclude)
	if err != nil {
		return err
	}
	if n > 0 {
		return apperror.Conflict("jersey number already taken in this team")
	}
	return nil
}

func (s *playerService) Create(ctx context.Context, p *models.Player) error {
	if err := s.ensureTeam(ctx, p.TeamID); err != nil {
		return err
	}
	if !models.IsValidPosition(p.Position) {
		return apperror.BadRequest("invalid position")
	}
	if err := s.jerseyAvailable(ctx, p.TeamID, p.JerseyNumber, nil); err != nil {
		return err
	}
	return s.players.Create(ctx, p)
}

func (s *playerService) Update(ctx context.Context, id uuid.UUID, patch *models.Player) (*models.Player, error) {
	existing, err := s.players.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("player not found")
		}
		return nil, err
	}
	teamID := existing.TeamID
	if patch.TeamID != uuid.Nil && patch.TeamID != existing.TeamID {
		if err := s.ensureTeam(ctx, patch.TeamID); err != nil {
			return nil, err
		}
		teamID = patch.TeamID
	}
	jersey := existing.JerseyNumber
	if patch.JerseyNumber != 0 {
		jersey = patch.JerseyNumber
	}
	if jersey != existing.JerseyNumber || teamID != existing.TeamID {
		if err := s.jerseyAvailable(ctx, teamID, jersey, &existing.ID); err != nil {
			return nil, err
		}
	}
	if patch.Name != "" {
		existing.Name = patch.Name
	}
	if patch.Height != 0 {
		existing.Height = patch.Height
	}
	if patch.Weight != 0 {
		existing.Weight = patch.Weight
	}
	if patch.Position != "" {
		if !models.IsValidPosition(patch.Position) {
			return nil, apperror.BadRequest("invalid position")
		}
		existing.Position = patch.Position
	}
	if patch.JerseyNumber != 0 {
		existing.JerseyNumber = patch.JerseyNumber
	}
	if patch.TeamID != uuid.Nil {
		existing.TeamID = patch.TeamID
	}
	if err := s.players.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *playerService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.players.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("player not found")
		}
		return err
	}
	return s.players.SoftDelete(ctx, id)
}

func (s *playerService) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]models.Player, error) {
	if _, err := s.teams.FindByIDUnscoped(ctx, teamID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("team not found")
		}
		return nil, err
	}
	return s.players.ListByTeam(ctx, teamID)
}
