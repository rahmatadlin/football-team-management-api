package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
)

type GoalRepository interface {
	ReplaceForMatch(ctx context.Context, matchID uuid.UUID, goals []models.Goal) error
	GoalsByMatch(ctx context.Context, matchID uuid.UUID) ([]models.Goal, error)
}

type goalRepo struct {
	db *gorm.DB
}

func NewGoalRepository(db *gorm.DB) GoalRepository {
	return &goalRepo{db: db}
}

func (r *goalRepo) ReplaceForMatch(ctx context.Context, matchID uuid.UUID, goals []models.Goal) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("match_id = ?", matchID).Delete(&models.Goal{}).Error; err != nil {
			return err
		}
		for i := range goals {
			goals[i].MatchID = matchID
			if err := tx.Create(&goals[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *goalRepo) GoalsByMatch(ctx context.Context, matchID uuid.UUID) ([]models.Goal, error) {
	var gs []models.Goal
	err := r.db.WithContext(ctx).Where("match_id = ?", matchID).Preload("Player").Find(&gs).Error
	return gs, err
}
