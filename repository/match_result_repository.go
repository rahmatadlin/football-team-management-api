package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
)

type MatchResultRepository interface {
	Upsert(ctx context.Context, r *models.MatchResult) error
	FindByMatchID(ctx context.Context, matchID uuid.UUID) (*models.MatchResult, error)
}

type matchResultRepo struct {
	db *gorm.DB
}

func NewMatchResultRepository(db *gorm.DB) MatchResultRepository {
	return &matchResultRepo{db: db}
}

func (r *matchResultRepo) Upsert(ctx context.Context, res *models.MatchResult) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.MatchResult
		err := tx.Where("match_id = ?", res.MatchID).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tx.Create(res).Error
		}
		if err != nil {
			return err
		}
		existing.HomeScore = res.HomeScore
		existing.AwayScore = res.AwayScore
		return tx.Save(&existing).Error
	})
}

func (r *matchResultRepo) FindByMatchID(ctx context.Context, matchID uuid.UUID) (*models.MatchResult, error) {
	var m models.MatchResult
	err := r.db.WithContext(ctx).Where("match_id = ?", matchID).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}
