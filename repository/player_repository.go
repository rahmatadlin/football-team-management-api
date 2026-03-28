package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
)

type PlayerRepository interface {
	Create(ctx context.Context, p *models.Player) error
	Update(ctx context.Context, p *models.Player) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Player, error)
	ListByTeam(ctx context.Context, teamID uuid.UUID) ([]models.Player, error)
	CountJerseyInTeam(ctx context.Context, teamID uuid.UUID, jersey int, excludePlayerID *uuid.UUID) (int64, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Player, error)
}

type playerRepo struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepo{db: db}
}

func (r *playerRepo) Create(ctx context.Context, p *models.Player) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *playerRepo) Update(ctx context.Context, p *models.Player) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *playerRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Player{}, "id = ?", id).Error
}

func (r *playerRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Player, error) {
	var p models.Player
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *playerRepo) ListByTeam(ctx context.Context, teamID uuid.UUID) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).Unscoped().Where("team_id = ?", teamID).Order("jersey_number ASC").Find(&players).Error
	return players, err
}

func (r *playerRepo) CountJerseyInTeam(ctx context.Context, teamID uuid.UUID, jersey int, excludePlayerID *uuid.UUID) (int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Player{}).Where("team_id = ? AND jersey_number = ?", teamID, jersey)
	if excludePlayerID != nil {
		q = q.Where("id <> ?", *excludePlayerID)
	}
	var count int64
	err := q.Count(&count).Error
	return count, err
}

func (r *playerRepo) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Player, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var players []models.Player
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&players).Error
	return players, err
}
