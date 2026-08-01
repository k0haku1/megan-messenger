package redis

import (
	"context"
	"errors"
	"megan-messenger/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrInvalidChallengeToken = errors.New("invalid challenge token")

type PasswordChallengeRepository struct {
	rdb       *redis.Client
	keyPrefix string
	ttl       time.Duration
}

func NewPasswordChallengeRepository(rdb *redis.Client, keyPrefix string, ttl time.Duration) repository.PasswordChallengeRepository {
	return &PasswordChallengeRepository{
		rdb:       rdb,
		keyPrefix: keyPrefix,
		ttl:       ttl,
	}
}

func (r *PasswordChallengeRepository) Issue(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	hash := hashToken(token)
	if err := r.rdb.Set(ctx, r.keyPrefix+hash, userID.String(), r.ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func (r *PasswordChallengeRepository) Consume(ctx context.Context, token string) (uuid.UUID, error) {
	hash := hashToken(token)
	userIDString, err := r.rdb.GetDel(ctx, r.keyPrefix+hash).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, ErrInvalidChallengeToken
		}
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}
