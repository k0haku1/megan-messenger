package repository

import (
	"context"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"

	"github.com/google/uuid"
)

type MessageRepository interface {
	ListMessages(ctx context.Context, conversationID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error)
	CreateMessage(ctx context.Context, msg model.Message) (model.Message, error)
}
