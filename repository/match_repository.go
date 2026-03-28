package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
)

type MatchRepository interface {
	Create(ctx context.Context, m *models.Match) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error)
	List(ctx context.Context) ([]models.Match, error)
	FindCompletedBefore(ctx context.Context, beforeDate string, beforeTime string, beforeID uuid.UUID) ([]models.Match, error)
}

type matchRepo struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepo{db: db}
}

func (r *matchRepo) Create(ctx context.Context, m *models.Match) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *matchRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Match, error) {
	var m models.Match
	err := r.db.WithContext(ctx).
		Preload("HomeTeam").
		Preload("AwayTeam").
		Preload("Result").
		Preload("Goals").
		Preload("Goals.Player").
		Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *matchRepo) List(ctx context.Context) ([]models.Match, error) {
	var list []models.Match
	err := r.db.WithContext(ctx).Unscoped().
		Preload("HomeTeam", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("AwayTeam", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Order("match_date ASC, match_time ASC, id ASC").
		Find(&list).Error
	return list, err
}

// FindCompletedBefore returns matches that have a result and occur strictly before the given match in sort order.
func (r *matchRepo) FindCompletedBefore(ctx context.Context, beforeDate string, beforeTime string, beforeID uuid.UUID) ([]models.Match, error) {
	var list []models.Match
	err := r.db.WithContext(ctx).
		Joins("JOIN match_results ON match_results.match_id = matches.id").
		Where(
			"(matches.match_date < ?) OR (matches.match_date = ? AND matches.match_time < ?) OR (matches.match_date = ? AND matches.match_time = ? AND matches.id < ?)",
			beforeDate, beforeDate, beforeTime, beforeDate, beforeTime, beforeID,
		).
		Preload("Result").
		Order("match_date ASC, match_time ASC, id ASC").
		Find(&list).Error
	return list, err
}
