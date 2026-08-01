package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Phone        string    `json:"phone"`
	Username     string    `json:"username,omitempty"`
	PasswordHash string    `json:"-"`
	AvatarURL    string    `json:"avatarUrl,omitempty"`
	CreatedAt    time.Time `json:"-"`
}

func NewUser(phone string) User {
	return User{
		ID:        uuid.Must(uuid.NewV7()),
		Phone:     phone,
		CreatedAt: time.Now(),
	}
}

func (u User) HasPassword() bool {
	return u.PasswordHash != ""
}

func (u User) OnboardingComplete() bool {
	return u.Username != ""
}

func (u *User) ComparePasswords(password string) error {
	if u.PasswordHash == "" {
		return errors.New("password not set")
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

func (u *User) ChangePassword(newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

type PhoneOTP struct {
	Phone     string
	CodeHash  string
	ExpiresAt time.Time
	Attempts  int
	SentAt    time.Time
	CreatedAt time.Time
}
