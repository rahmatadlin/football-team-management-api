package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	TeamID       uuid.UUID `json:"team_id" gorm:"type:char(36);not null;index"`
	Team         Team      `json:"-" gorm:"foreignKey:TeamID"`
	Name         string    `json:"name" gorm:"size:255;not null"`
	Height       float64   `json:"height" gorm:"type:decimal(5,2);not null"`
	Weight       float64   `json:"weight" gorm:"type:decimal(5,2);not null"`
	Position     Position  `json:"position" gorm:"type:varchar(32);not null"`
	JerseyNumber int       `json:"jersey_number" gorm:"not null"`
	SoftDeleteModel
	Timestamps
}

func (p *Player) BeforeCreate(tx *gorm.DB) error {
	return BeforeCreateUUID(&p.ID)
}

func (Player) TableName() string {
	return "players"
}
