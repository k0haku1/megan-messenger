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
)

type FolderRepository struct {
	queries db.Querier
}

func NewFolderRepository(queries db.Querier) repository.FolderRepository {
	return &FolderRepository{queries: queries}
}

func mapChatFolder(row db.ChatFolder) model.ChatFolder {
	return model.ChatFolder{
		ID:        row.ID,
		UserID:    row.UserID,
		Name:      row.Name,
		Icon:      textOrEmpty(row.Icon),
		Position:  int(row.Position),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func (r *FolderRepository) Create(ctx context.Context, folder model.ChatFolder) (model.ChatFolder, error) {
	row, err := r.queries.CreateChatFolder(ctx, db.CreateChatFolderParams{
		ID:        folder.ID,
		UserID:    folder.UserID,
		Name:      folder.Name,
		Icon:      textFromString(folder.Icon),
		Position:  int32(folder.Position),
		CreatedAt: timestampFromTime(folder.CreatedAt),
		UpdatedAt: timestampFromTime(folder.UpdatedAt),
	})
	if err != nil {
		return model.ChatFolder{}, err
	}
	return mapChatFolder(row), nil
}

func (r *FolderRepository) GetByID(ctx context.Context, id uuid.UUID) (model.ChatFolder, error) {
	row, err := r.queries.GetChatFolder(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChatFolder{}, repository.ErrFolderNotFound
		}
		return model.ChatFolder{}, err
	}
	return mapChatFolder(row), nil
}

func (r *FolderRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.ChatFolder, error) {
	rows, err := r.queries.ListChatFolders(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]model.ChatFolder, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapChatFolder(row))
	}
	return items, nil
}

func (r *FolderRepository) Update(ctx context.Context, folder model.ChatFolder) (model.ChatFolder, error) {
	row, err := r.queries.UpdateChatFolder(ctx, db.UpdateChatFolderParams{
		ID:        folder.ID,
		UserID:    folder.UserID,
		Name:      folder.Name,
		Icon:      textFromString(folder.Icon),
		UpdatedAt: timestampFromTime(time.Now()),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChatFolder{}, repository.ErrFolderNotFound
		}
		return model.ChatFolder{}, err
	}
	return mapChatFolder(row), nil
}

func (r *FolderRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.DeleteChatFolder(ctx, db.DeleteChatFolderParams{ID: id, UserID: userID})
}

func (r *FolderRepository) UpdatePosition(ctx context.Context, id uuid.UUID, position int, userID uuid.UUID) error {
	return r.queries.UpdateChatFolderPosition(ctx, db.UpdateChatFolderPositionParams{
		ID:       id,
		Position: int32(position),
		UserID:   userID,
	})
}

func (r *FolderRepository) MaxPosition(ctx context.Context, userID uuid.UUID) (int, error) {
	v, err := r.queries.MaxChatFolderPosition(ctx, userID)
	return int(v), err
}

func (r *FolderRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	v, err := r.queries.CountUserChatFolders(ctx, userID)
	return int(v), err
}

func (r *FolderRepository) AddItem(ctx context.Context, folderID, conversationID uuid.UUID) error {
	return r.queries.AddChatFolderItem(ctx, db.AddChatFolderItemParams{
		FolderID:       folderID,
		ConversationID: conversationID,
		AddedAt:        timestampFromTime(time.Now()),
	})
}

func (r *FolderRepository) RemoveItem(ctx context.Context, folderID, conversationID uuid.UUID) error {
	return r.queries.RemoveChatFolderItem(ctx, db.RemoveChatFolderItemParams{
		FolderID:       folderID,
		ConversationID: conversationID,
	})
}

func (r *FolderRepository) ItemExists(ctx context.Context, folderID, conversationID uuid.UUID) (bool, error) {
	return r.queries.ChatFolderItemExists(ctx, db.ChatFolderItemExistsParams{
		FolderID:       folderID,
		ConversationID: conversationID,
	})
}

func (r *FolderRepository) CountItems(ctx context.Context, folderID uuid.UUID) (int, error) {
	v, err := r.queries.CountChatFolderItems(ctx, folderID)
	return int(v), err
}

func (r *FolderRepository) GetFolderIDsByConversations(ctx context.Context, userID uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	rows, err := r.queries.GetFolderIDsByConversations(ctx, userID)
	if err != nil {
		return nil, err
	}
	m := make(map[uuid.UUID][]uuid.UUID, len(rows))
	for _, row := range rows {
		m[row.ConversationID] = append(m[row.ConversationID], row.FolderID)
	}
	return m, nil
}
