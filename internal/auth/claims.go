package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Phone              string `json:"phone"`
	Username           string `json:"username,omitempty"`
	OnboardingComplete bool   `json:"onboardingComplete"`

	jwt.RegisteredClaims
}
