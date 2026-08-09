package postgres

import (
	"context"
	"errors"
	db "megan-messenger/internal/db/postgres/sqlc"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"
	"megan-messenger/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageRepository struct {
	queries *db.Queries
}

func NewMessageRepository(queries *db.Queries) *MessageRepository {
	return &MessageRepository{queries: queries}
}

func mapMessage(message db.Message, sender model.MessageSender) model.Message {
	return model.Message{
		ID:              message.ID,
		ConversationID:  message.ConversationID,
		Content:         message.Content,
		ReplyToID:       uuidPtr(message.ReplyToID),
		ForwardedFromID: uuidPtr(message.ForwardedFromID),
		DeletedAt:       timePtr(message.DeletedAt),
		Sender:          sender,
		CreatedAt:       message.CreatedAt.Time,
	}
}

func mapMessages(rows []db.GetMessagesPagingRow) []model.Message {
	result := make([]model.Message, 0, len(rows))
	for _, r := range rows {
		result = append(result, model.Message{
			ID:              r.ID,
			ConversationID:  r.ConversationID,
			Content:         r.Content,
			ReplyToID:       uuidPtr(r.ReplyToID),
			ForwardedFromID: uuidPtr(r.ForwardedFromID),
			DeletedAt:       timePtr(r.DeletedAt),
			CreatedAt:       r.CreatedAt.Time,
			Sender: model.MessageSender{
				ID:        r.SenderID,
				Username:  textOrEmpty(r.Username),
				AvatarURL: textOrEmpty(r.AvatarUrl),
			},
		})
	}
	return result
}

func (r *MessageRepository) CreateMessage(ctx context.Context, msg model.Message) (model.Message, error) {
	createdMessage, err := r.queries.CreateMessage(ctx, db.CreateMessageParams{
		ID:              msg.ID,
		ConversationID:  msg.ConversationID,
		SenderID:        msg.Sender.ID,
		Content:         msg.Content,
		ReplyToID:       uuidFromPtr(msg.ReplyToID),
		ForwardedFromID: uuidFromPtr(msg.ForwardedFromID),
		CreatedAt:       timestampFromTime(msg.CreatedAt),
	})
	if err != nil {
		return model.Message{}, err
	}

	return mapMessage(createdMessage, msg.Sender), err
}

func (r *MessageRepository) ListMessages(ctx context.Context, conversationID, userID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error) {
	params := db.GetMessagesPagingParams{
		ConversationID: conversationID,
		UserID:         userID,
		Limit:          int32(limit),
	}

	if cursor != nil {
		params.CursorID = cursor.ID
		params.CursorCreatedAt = pgtype.Timestamptz{
			Time:  cursor.CreatedAt,
			Valid: true,
		}
	}

	messages, err := r.queries.GetMessagesPaging(ctx, params)
	if err != nil {
		return nil, err
	}

	result := mapMessages(messages)
	for i := range result {
		reactions, err := r.reactions(ctx, result[i].ID, userID)
		if err != nil {
			return nil, err
		}
		result[i].Reactions = reactions
	}
	return result, nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id uuid.UUID) (model.Message, error) {
	row, err := r.queries.GetMessageByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Message{}, repository.ErrMessageNotFound
		}
		return model.Message{}, err
	}

	return model.Message{
		ID:              row.ID,
		ConversationID:  row.ConversationID,
		Content:         row.Content,
		ReplyToID:       uuidPtr(row.ReplyToID),
		ForwardedFromID: uuidPtr(row.ForwardedFromID),
		DeletedAt:       timePtr(row.DeletedAt),
		CreatedAt:       row.CreatedAt.Time,
		Sender: model.MessageSender{
			ID:        row.SenderID,
			Username:  textOrEmpty(row.Username),
			AvatarURL: textOrEmpty(row.AvatarUrl),
		},
	}, nil
}

func (r *MessageRepository) GetByIDForUser(ctx context.Context, id, userID uuid.UUID) (model.Message, error) {
	message, err := r.GetByID(ctx, id)
	if err != nil {
		return model.Message{}, err
	}
	reactions, err := r.reactions(ctx, id, userID)
	if err != nil {
		return model.Message{}, err
	}
	message.Reactions = reactions
	return message, nil
}

func (r *MessageRepository) reactions(ctx context.Context, messageID, userID uuid.UUID) ([]model.MessageReaction, error) {
	rows, err := r.queries.GetMessageReactions(ctx, db.GetMessageReactionsParams{MessageID: messageID, UserID: userID})
	if err != nil {
		return nil, err
	}
	result := make([]model.MessageReaction, 0, len(rows))
	byEmoji := make(map[string]int, len(rows))
	for _, row := range rows {
		byEmoji[row.Emoji] = len(result)
		result = append(result, model.MessageReaction{Emoji: row.Emoji, Count: row.Count, ReactedByMe: row.ReactedByMe, ReactedBy: []model.MessageSender{}})
	}
	users, err := r.queries.GetMessageReactionUsers(ctx, messageID)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if index, ok := byEmoji[user.Emoji]; ok {
			result[index].ReactedBy = append(result[index].ReactedBy, model.MessageSender{ID: user.ID, Username: textOrEmpty(user.Username), AvatarURL: textOrEmpty(user.AvatarUrl)})
		}
	}
	return result, nil
}

func (r *MessageRepository) HideForUser(ctx context.Context, messageID, userID uuid.UUID) error {
	return r.queries.HideMessageForUser(ctx, db.HideMessageForUserParams{MessageID: messageID, UserID: userID})
}

func (r *MessageRepository) DeleteForEveryone(ctx context.Context, messageID, userID uuid.UUID) (model.Message, error) {
	message, err := r.queries.DeleteMessageForEveryone(ctx, db.DeleteMessageForEveryoneParams{ID: messageID, SenderID: userID})
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Message{}, repository.ErrMessageNotFound
	}
	if err != nil {
		return model.Message{}, err
	}
	return mapMessage(message, model.MessageSender{ID: userID}), nil
}

func (r *MessageRepository) AddReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	return r.queries.AddMessageReaction(ctx, db.AddMessageReactionParams{MessageID: messageID, UserID: userID, Emoji: emoji})
}

func (r *MessageRepository) RemoveReaction(ctx context.Context, messageID, userID uuid.UUID, emoji string) error {
	return r.queries.RemoveMessageReaction(ctx, db.RemoveMessageReactionParams{MessageID: messageID, UserID: userID, Emoji: emoji})
}
