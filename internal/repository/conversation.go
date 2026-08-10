package repository

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

var ErrConversationNotFound = errors.New("conversation not found")

type ConversationRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Conversation, error)
	GetBySlug(ctx context.Context, slug string) (model.Conversation, error)
	Create(ctx context.Context, conversation model.Conversation) (model.Conversation, error)
	AddMember(ctx context.Context, conversationID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, conversationID, userID uuid.UUID) error
	IsMember(ctx context.Context, conversationID, userID uuid.UUID) (bool, error)
	FindDM(ctx context.Context, userLow, userHigh uuid.UUID) (uuid.UUID, error)
	CreateDMPair(ctx context.Context, userLow, userHigh, conversationID uuid.UUID) error
	AdvanceLastReadAt(ctx context.Context, conversationID, userID uuid.UUID, readAt time.Time) (time.Time, bool, error)
	GetLastReadAt(ctx context.Context, conversationID, userID uuid.UUID) (*time.Time, error)
	GetOthersReadWatermark(ctx context.Context, conversationID, viewerID uuid.UUID) (*time.Time, error)
	ListMemberIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error)
}
