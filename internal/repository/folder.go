package repository

import (
	"context"
	"errors"
	"megan-messenger/internal/model"

	"github.com/google/uuid"
)

var ErrFolderNotFound = errors.New("folder not found")

type FolderRepository interface {
	Create(ctx context.Context, folder model.ChatFolder) (model.ChatFolder, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.ChatFolder, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.ChatFolder, error)
	Update(ctx context.Context, folder model.ChatFolder) (model.ChatFolder, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	UpdatePosition(ctx context.Context, id uuid.UUID, position int, userID uuid.UUID) error
	MaxPosition(ctx context.Context, userID uuid.UUID) (int, error)
	CountByUser(ctx context.Context, userID uuid.UUID) (int, error)

	AddItem(ctx context.Context, folderID, conversationID uuid.UUID) error
	RemoveItem(ctx context.Context, folderID, conversationID uuid.UUID) error
	ItemExists(ctx context.Context, folderID, conversationID uuid.UUID) (bool, error)
	CountItems(ctx context.Context, folderID uuid.UUID) (int, error)
	GetFolderIDsByConversations(ctx context.Context, userID uuid.UUID) (map[uuid.UUID][]uuid.UUID, error)
}
