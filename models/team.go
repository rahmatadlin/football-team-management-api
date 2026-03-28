package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID           uuid.UUID      `json:"id" gorm:"type:char(36);primaryKey"`
	Name         string         `json:"name" gorm:"size:255;not null"`
	LogoURL      string         `json:"logo_url" gorm:"size:512"`
	FoundedYear  int            `json:"founded_year" gorm:"not null"`
	Address      string         `json:"address" gorm:"size:512"`
	City         string         `json:"city" gorm:"size:128"`
	Players      []Player       `json:"players,omitempty" gorm:"foreignKey:TeamID"`
	SoftDeleteModel
	Timestamps
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	return BeforeCreateUUID(&t.ID)
}

func (Team) TableName() string {
	return "teams"
}
