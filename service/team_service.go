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

type TeamUpdateInput struct {
	Name        *string
	LogoURL     *string
	FoundedYear *int
	Address     *string
	City        *string
}

type TeamService interface {
	Create(ctx context.Context, t *models.Team) error
	Update(ctx context.Context, id uuid.UUID, patch *TeamUpdateInput) (*models.Team, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*models.Team, error)
	List(ctx context.Context) ([]models.Team, error)
}

type teamService struct {
	teams repository.TeamRepository
}

func NewTeamService(teams repository.TeamRepository) TeamService {
	return &teamService{teams: teams}
}

func (s *teamService) Create(ctx context.Context, t *models.Team) error {
	if t.Name == "" {
		return apperror.BadRequest("name is required")
	}
	return s.teams.Create(ctx, t)
}

func (s *teamService) Update(ctx context.Context, id uuid.UUID, patch *TeamUpdateInput) (*models.Team, error) {
	existing, err := s.teams.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("team not found")
		}
		return nil, err
	}
	if patch == nil {
		return existing, nil
	}
	if patch.Name != nil {
		existing.Name = *patch.Name
	}
	if patch.LogoURL != nil {
		existing.LogoURL = *patch.LogoURL
	}
	if patch.FoundedYear != nil {
		existing.FoundedYear = *patch.FoundedYear
	}
	if patch.Address != nil {
		existing.Address = *patch.Address
	}
	if patch.City != nil {
		existing.City = *patch.City
	}
	if err := s.teams.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *teamService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.teams.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("team not found")
		}
		return err
	}
	return s.teams.SoftDelete(ctx, id)
}

func (s *teamService) Get(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	t, err := s.teams.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("team not found")
		}
		return nil, err
	}
	return t, nil
}

func (s *teamService) List(ctx context.Context) ([]models.Team, error) {
	return s.teams.List(ctx)
}
