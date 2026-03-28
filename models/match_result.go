package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MatchResult struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	MatchID   uuid.UUID `json:"match_id" gorm:"type:char(36);not null;uniqueIndex"`
	Match     Match     `json:"-" gorm:"foreignKey:MatchID"`
	HomeScore int       `json:"home_score" gorm:"not null"`
	AwayScore int       `json:"away_score" gorm:"not null"`
	Timestamps
}

func (r *MatchResult) BeforeCreate(tx *gorm.DB) error {
	return BeforeCreateUUID(&r.ID)
}

func (MatchResult) TableName() string {
	return "match_results"
}
