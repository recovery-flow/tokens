package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Role     string    `json:"Role"`
	DeviceID uuid.UUID `json:"device_id"`
}

func (m *TokenManager) GenerateJWT(
	userID uuid.UUID,
	deviceID uuid.UUID,
	role string,
	tlt time.Duration,
	sk string,
) (string, error) {
	if !IsSupportedRole(role) {
		return "", errors.New("unsupported role")
	}

	expirationTime := time.Now().Add(tlt * time.Second)
	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Role:     role,
		DeviceID: deviceID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

type ServiceClaims struct {
	jwt.RegisteredClaims
	Service string `json:"service"`
}

func (m *TokenManager) GenerateServiceJWT(serviceName string, tlt time.Duration, sk string) (string, error) {
	expirationTime := time.Now().Add(tlt)

	claims := &ServiceClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   serviceName,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
		Service: serviceName,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}
