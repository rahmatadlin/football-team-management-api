package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"football-team-management-api/config"
	"football-team-management-api/models"
	"football-team-management-api/repository"
	"football-team-management-api/utils/apperror"
	jwtutil "football-team-management-api/utils/jwt"
	"football-team-management-api/utils/password"
)

type AuthService interface {
	Login(ctx context.Context, email, plainPassword string) (token string, admin *models.Admin, err error)
}

type authService struct {
	cfg    *config.Config
	admins repository.AdminRepository
}

func NewAuthService(cfg *config.Config, admins repository.AdminRepository) AuthService {
	return &authService{cfg: cfg, admins: admins}
}

func (s *authService) Login(ctx context.Context, email, plainPassword string) (string, *models.Admin, error) {
	a, err := s.admins.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, apperror.Unauthorized("invalid email or password")
		}
		return "", nil, err
	}
	if !password.Verify(a.PasswordHash, plainPassword) {
		return "", nil, apperror.Unauthorized("invalid email or password")
	}
	exp := time.Duration(s.cfg.JWTExpiryHrs) * time.Hour
	token, err := jwtutil.Sign([]byte(s.cfg.JWTSecret), a.ID, a.Email, exp)
	if err != nil {
		return "", nil, apperror.Internal("could not issue token")
	}
	return token, a, nil
}
