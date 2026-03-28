package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
)

type AdminRepository interface {
	FindByEmail(ctx context.Context, email string) (*models.Admin, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.Admin, error)
	Create(ctx context.Context, a *models.Admin) error
}

type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepo{db: db}
}

func (r *adminRepo) FindByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var a models.Admin
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *adminRepo) FindByID(ctx context.Context, id uuid.UUID) (*models.Admin, error) {
	var a models.Admin
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *adminRepo) Create(ctx context.Context, a *models.Admin) error {
	return r.db.WithContext(ctx).Create(a).Error
}
