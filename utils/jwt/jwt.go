package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	AdminID uuid.UUID `json:"admin_id"`
	Email   string    `json:"email"`
	jwtlib.RegisteredClaims
}

func Sign(secret []byte, adminID uuid.UUID, email string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		AdminID: adminID,
		Email:   email,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			Subject:   adminID.String(),
		},
	}
	t := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func Parse(secret []byte, tokenStr string) (*Claims, error) {
	t, err := jwtlib.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
