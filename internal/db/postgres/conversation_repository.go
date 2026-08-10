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
		ID:          row.ID,
		Type:        model.ConversationType(row.Type),
		Title:       textOrEmpty(row.Title),
		Slug:        textOrEmpty(row.Slug),
		UnreadCount: int(row.UnreadCount),
		OthersReadAt: timePtr(row.OthersReadAt),
		CreatedAt:   row.CreatedAt.Time,
	}

	if row.PeerID.Valid {
		conv.Peer = &model.ConversationPeer{
			ID:        uuid.UUID(row.PeerID.Bytes),
			Username:  textOrEmpty(row.PeerUsername),
			AvatarURL: textOrEmpty(row.PeerAvatarUrl),
		}
	}

	if row.LastMessageID.Valid {
		senderID := uuid.UUID{}
		if row.LastMessageSenderID.Valid {
			senderID = uuid.UUID(row.LastMessageSenderID.Bytes)
		}
		conv.LastMessage = &model.ConversationPreview{
			ID:             uuid.UUID(row.LastMessageID.Bytes),
			Content:        textOrEmpty(row.LastMessageContent),
			SenderID:       senderID,
			SenderUsername: textOrEmpty(row.LastMessageSenderUsername),
			AttachmentKind: row.LastMessageAttachmentKind,
			CreatedAt:      row.LastMessageCreatedAt.Time,
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

func (r *ConversationRepository) AdvanceLastReadAt(
	ctx context.Context,
	conversationID, userID uuid.UUID,
	readAt time.Time,
) (time.Time, bool, error) {
	updated, err := r.queries.AdvanceMemberLastReadAt(ctx, db.AdvanceMemberLastReadAtParams{
		LastReadAt:     timestampFromTime(readAt),
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err == nil {
		return updated.Time, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, false, err
	}

	current, err := r.queries.GetMemberLastReadAt(ctx, db.GetMemberLastReadAtParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, false, repository.ErrConversationNotFound
		}
		return time.Time{}, false, err
	}
	if !current.Valid {
		return time.Time{}, false, nil
	}
	return current.Time, false, nil
}

func (r *ConversationRepository) GetLastReadAt(
	ctx context.Context,
	conversationID, userID uuid.UUID,
) (*time.Time, error) {
	value, err := r.queries.GetMemberLastReadAt(ctx, db.GetMemberLastReadAtParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrConversationNotFound
		}
		return nil, err
	}
	return timePtr(value), nil
}

func (r *ConversationRepository) GetOthersReadWatermark(
	ctx context.Context,
	conversationID, viewerID uuid.UUID,
) (*time.Time, error) {
	value, err := r.queries.GetOthersReadWatermark(ctx, db.GetOthersReadWatermarkParams{
		ConversationID: conversationID,
		ViewerID:       viewerID,
	})
	if err != nil {
		return nil, err
	}
	return timePtr(value), nil
}

func (r *ConversationRepository) ListMemberIDs(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	return r.queries.ListConversationMemberIDs(ctx, conversationID)
}
