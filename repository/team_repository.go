package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
)

type TeamRepository interface {
	Create(ctx context.Context, t *models.Team) error
	Update(ctx context.Context, t *models.Team) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Team, error)
	FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Team, error)
	List(ctx context.Context) ([]models.Team, error)
}

type teamRepo struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(ctx context.Context, t *models.Team) error {
	return r.db.WithContext(ctx).Create(t).Error
}

func (r *teamRepo) Update(ctx context.Context, t *models.Team) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *teamRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Team{}, "id = ?", id).Error
}

func (r *teamRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	var t models.Team
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *teamRepo) FindByIDUnscoped(ctx context.Context, id uuid.UUID) (*models.Team, error) {
	var t models.Team
	err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *teamRepo) List(ctx context.Context) ([]models.Team, error) {
	var teams []models.Team
	err := r.db.WithContext(ctx).Unscoped().Order("created_at DESC").Find(&teams).Error
	return teams, err
}
