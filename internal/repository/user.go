package repository

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

var ErrUniqueAlreadyExists = errors.New("value already exists for unique field")
var ErrUserNotFound = errors.New("user not found")
var ErrOTPNotFound = errors.New("otp not found")

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByPhone(ctx context.Context, phone string) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Create(ctx context.Context, u model.User) (model.User, error)
	SetUsername(ctx context.Context, id uuid.UUID, username string) error
	UpdatePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error
	ClearPassword(ctx context.Context, id uuid.UUID) error
	ChangeAvatar(ctx context.Context, id uuid.UUID, url string) error
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
	SearchByUsernamePrefix(ctx context.Context, prefix string, excludeUserID uuid.UUID, limit int32) ([]model.UserSearchResult, error)
	UpdateUsername(ctx context.Context, id uuid.UUID, username string) error
	UpdatePrivacy(ctx context.Context, id uuid.UUID, searchable bool, policy model.DMPolicy) error

	UpsertPhoneOTP(ctx context.Context, phone, codeHash string, expiresAt, sentAt time.Time) error
	GetPhoneOTP(ctx context.Context, phone string) (model.PhoneOTP, error)
	IncrementPhoneOTPAttempts(ctx context.Context, phone string) error
	DeletePhoneOTP(ctx context.Context, phone string) error
}

type PasswordChallengeRepository interface {
	Issue(ctx context.Context, userID uuid.UUID) (string, error)
	Consume(ctx context.Context, token string) (uuid.UUID, error)
}
