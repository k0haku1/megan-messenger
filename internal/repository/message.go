package repository

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"

	"github.com/google/uuid"
)

var ErrMessageNotFound = errors.New("message not found")

type MessageRepository interface {
	ListMessages(ctx context.Context, conversationID, userID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error)
	CreateMessage(ctx context.Context, msg model.Message) (model.Message, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Message, error)
	GetByIDForUser(ctx context.Context, id, userID uuid.UUID) (model.Message, error)
	HideForUser(ctx context.Context, messageID, userID uuid.UUID) error
	DeleteForEveryone(ctx context.Context, messageID, userID uuid.UUID) (model.Message, error)
	AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error
	RemoveReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error
}
