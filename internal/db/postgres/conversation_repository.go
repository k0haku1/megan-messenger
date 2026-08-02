package postgres

import (
	"context"
	"errors"
	db "megan-messenger/internal/db/postgres/sqlc"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ConversationRepository struct {
	queries db.Querier
}

func NewConversationRepository(queries db.Querier) repository.ConversationRepository {
	return &ConversationRepository{queries: queries}
}

func mapConversationRow(row db.Conversation) model.Conversation {
	return model.Conversation{
		ID:        row.ID,
		Type:      model.ConversationType(row.Type),
		Title:     textOrEmpty(row.Title),
		Slug:      textOrEmpty(row.Slug),
		CreatedAt: row.CreatedAt.Time,
	}
}

func mapListedConversation(row db.GetUserConversationsRow) model.Conversation {
	conv := model.Conversation{
		ID:        row.ID,
		Type:      model.ConversationType(row.Type),
		Title:     textOrEmpty(row.Title),
		Slug:      textOrEmpty(row.Slug),
		CreatedAt: row.CreatedAt.Time,
	}

	if row.PeerID.Valid {
		conv.Peer = &model.ConversationPeer{
			ID:        uuid.UUID(row.PeerID.Bytes),
			Username:  textOrEmpty(row.PeerUsername),
			AvatarURL: textOrEmpty(row.PeerAvatarUrl),
		}
	}

	return conv
}

func mapConversations(rows []db.GetUserConversationsRow) []model.Conversation {
	result := make([]model.Conversation, len(rows))
	for i, row := range rows {
		result[i] = mapListedConversation(row)
	}
	return result
}

func (r *ConversationRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	rows, err := r.queries.GetUserConversations(ctx, userID)
	if err != nil {
		return nil, err
	}
	return mapConversations(rows), nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Conversation, error) {
	row, err := r.queries.GetConversation(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Conversation{}, repository.ErrConversationNotFound
		}
		return model.Conversation{}, err
	}
	return mapConversationRow(row), nil
}

func (r *ConversationRepository) GetBySlug(ctx context.Context, slug string) (model.Conversation, error) {
	row, err := r.queries.GetConversationBySlug(ctx, textFromString(slug))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Conversation{}, repository.ErrConversationNotFound
		}
		return model.Conversation{}, err
	}
	return mapConversationRow(row), nil
}

func (r *ConversationRepository) Create(ctx context.Context, conversation model.Conversation) (model.Conversation, error) {
	created, err := r.queries.CreateConversation(ctx, db.CreateConversationParams{
		ID:        conversation.ID,
		Type:      db.ConversationType(conversation.Type),
		Title:     textFromString(conversation.Title),
		Slug:      textFromString(conversation.Slug),
		CreatedAt: timestampFromTime(conversation.CreatedAt),
	})
	if err != nil {
		return model.Conversation{}, err
	}
	return mapConversationRow(created), nil
}

func (r *ConversationRepository) AddMember(ctx context.Context, conversationID, userID uuid.UUID) error {
	member := model.NewConversationMember(userID, conversationID)
	return r.queries.AddConversationMember(ctx, db.AddConversationMemberParams{
		ID:             member.ID,
		ConversationID: member.ConversationID,
		UserID:         member.UserID,
		JoinedAt:       timestampFromTime(member.JoinedAt),
	})
}

func (r *ConversationRepository) RemoveMember(ctx context.Context, conversationID, userID uuid.UUID) error {
	return r.queries.RemoveConversationMember(ctx, db.RemoveConversationMemberParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
}

func (r *ConversationRepository) IsMember(ctx context.Context, conversationID, userID uuid.UUID) (bool, error) {
	return r.queries.IsConversationMember(ctx, db.IsConversationMemberParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
}

func (r *ConversationRepository) FindDM(ctx context.Context, userLow, userHigh uuid.UUID) (uuid.UUID, error) {
	return r.queries.FindDMConversation(ctx, db.FindDMConversationParams{
		UserLow:  userLow,
		UserHigh: userHigh,
	})
}

func (r *ConversationRepository) CreateDMPair(ctx context.Context, userLow, userHigh, conversationID uuid.UUID) error {
	return r.queries.CreateDMPair(ctx, db.CreateDMPairParams{
		UserLow:        userLow,
		UserHigh:       userHigh,
		ConversationID: conversationID,
	})
}
