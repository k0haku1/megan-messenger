package repository

import (
	"context"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"

	"github.com/google/uuid"
)

type AttachmentRepository interface {
	Create(ctx context.Context, att model.MessageAttachment) (model.MessageAttachment, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.MessageAttachment, error)
	ListByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) ([]model.MessageAttachment, error)
	AttachToMessage(ctx context.Context, messageID, uploaderID, conversationID uuid.UUID, attachmentIDs []uuid.UUID) error
	CountPending(ctx context.Context, uploaderID, conversationID uuid.UUID, attachmentIDs []uuid.UUID) (int64, error)
	ListConversationMedia(ctx context.Context, conversationID uuid.UUID, kind string, limit int, cursor *pagination.Cursor) ([]model.MessageAttachment, error)
}
