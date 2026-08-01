package postgres

import (
	"context"
	"errors"
	db "megan-messenger/internal/db/postgres/sqlc"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	queries db.Querier
}

func NewUserRepository(queries db.Querier) repository.UserRepository {
	return &UserRepository{queries}
}

func mapUser(user db.User) model.User {
	return model.User{
		ID:           user.ID,
		Phone:        user.Phone,
		Username:     textOrEmpty(user.Username),
		PasswordHash: textOrEmpty(user.PasswordHash),
		AvatarURL:    textOrEmpty(user.AvatarUrl),
		CreatedAt:    user.CreatedAt.Time,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	u, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repository.ErrUserNotFound
		}
		return model.User{}, err
	}
	return mapUser(u), nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (model.User, error) {
	u, err := r.queries.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repository.ErrUserNotFound
		}
		return model.User{}, err
	}
	return mapUser(u), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	u, err := r.queries.GetUserByUsername(ctx, pgtype.Text{String: username, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, repository.ErrUserNotFound
		}
		return model.User{}, err
	}
	return mapUser(u), nil
}

func (r *UserRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	created, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:        u.ID,
		Phone:     u.Phone,
		CreatedAt: timestampFromTime(u.CreatedAt),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.User{}, repository.ErrUniqueAlreadyExists
		}
		return model.User{}, err
	}
	return mapUser(created), nil
}

func (r *UserRepository) SetUsername(ctx context.Context, id uuid.UUID, username string) error {
	err := r.queries.SetUsername(ctx, db.SetUsernameParams{
		Username: pgtype.Text{String: username, Valid: true},
		ID:       id,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrUniqueAlreadyExists
		}
	}
	return err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error {
	return r.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           id,
		PasswordHash: pgtype.Text{String: newPasswordHash, Valid: true},
	})
}

func (r *UserRepository) ClearPassword(ctx context.Context, id uuid.UUID) error {
	return r.queries.ClearUserPassword(ctx, id)
}

func (r *UserRepository) ChangeAvatar(ctx context.Context, id uuid.UUID, url string) error {
	return r.queries.UpdateUserAvatar(ctx, db.UpdateUserAvatarParams{
		AvatarUrl: pgtype.Text{String: url, Valid: true},
		ID:        id,
	})
}

func (r *UserRepository) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	return r.queries.UserWithUsernameExists(ctx, pgtype.Text{String: username, Valid: true})
}

func (r *UserRepository) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	return r.queries.UserWithPhoneExists(ctx, phone)
}

func (r *UserRepository) UpsertPhoneOTP(ctx context.Context, phone, codeHash string, expiresAt, sentAt time.Time) error {
	now := time.Now()
	return r.queries.UpsertPhoneOTP(ctx, db.UpsertPhoneOTPParams{
		Phone:     phone,
		CodeHash:  codeHash,
		ExpiresAt: timestampFromTime(expiresAt),
		Attempts:  0,
		SentAt:    timestampFromTime(sentAt),
		CreatedAt: timestampFromTime(now),
	})
}

func (r *UserRepository) GetPhoneOTP(ctx context.Context, phone string) (model.PhoneOTP, error) {
	otp, err := r.queries.GetPhoneOTP(ctx, phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PhoneOTP{}, repository.ErrOTPNotFound
		}
		return model.PhoneOTP{}, err
	}
	return model.PhoneOTP{
		Phone:     otp.Phone,
		CodeHash:  otp.CodeHash,
		ExpiresAt: otp.ExpiresAt.Time,
		Attempts:  int(otp.Attempts),
		SentAt:    otp.SentAt.Time,
		CreatedAt: otp.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) IncrementPhoneOTPAttempts(ctx context.Context, phone string) error {
	return r.queries.IncrementPhoneOTPAttempts(ctx, phone)
}

func (r *UserRepository) DeletePhoneOTP(ctx context.Context, phone string) error {
	return r.queries.DeletePhoneOTP(ctx, phone)
}
