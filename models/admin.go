package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `json:"-" gorm:"column:password_hash;size:255;not null"`
	Timestamps
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	return BeforeCreateUUID(&a.ID)
}

func (Admin) TableName() string {
	return "admins"
}
