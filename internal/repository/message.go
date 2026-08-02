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
	ListMessages(ctx context.Context, conversationID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error)
	CreateMessage(ctx context.Context, msg model.Message) (model.Message, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Message, error)
}
