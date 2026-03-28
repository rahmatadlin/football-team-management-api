package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func BeforeCreateUUID(id *uuid.UUID) error {
	if *id == uuid.Nil {
		*id = uuid.New()
	}
	return nil
}

type SoftDeleteModel struct {
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
