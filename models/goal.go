package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Goal struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	MatchID  uuid.UUID `json:"match_id" gorm:"type:char(36);not null;index"`
	Match    Match     `json:"-" gorm:"foreignKey:MatchID"`
	PlayerID uuid.UUID `json:"player_id" gorm:"type:char(36);not null;index"`
	Player   Player    `json:"player,omitempty" gorm:"foreignKey:PlayerID"`
	GoalTime int       `json:"goal_time" gorm:"not null"` // minute
	Timestamps
}

func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	return BeforeCreateUUID(&g.ID)
}

func (Goal) TableName() string {
	return "goals"
}
