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
		ID:             message.ID,
		ConversationID: message.ConversationID,
		Content:        message.Content,
		Sender:         sender,
		CreatedAt:      message.CreatedAt.Time,
	}
}

func mapMessages(rows []db.GetMessagesPagingRow) []model.Message {
	result := make([]model.Message, 0, len(rows))
	for _, r := range rows {
		result = append(result, model.Message{
			ID:             r.ID,
			ConversationID: r.ConversationID,
			Content:        r.Content,
			CreatedAt:      r.CreatedAt.Time,
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
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		SenderID:       msg.Sender.ID,
		Content:        msg.Content,
		CreatedAt:      timestampFromTime(msg.CreatedAt),
	})
	if err != nil {
		return model.Message{}, err
	}

	return mapMessage(createdMessage, msg.Sender), err
}

func (r *MessageRepository) ListMessages(ctx context.Context, conversationID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error) {
	params := db.GetMessagesPagingParams{
		ConversationID: conversationID,
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

	return mapMessages(messages), nil
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
		ID:             row.ID,
		ConversationID: row.ConversationID,
		Content:        row.Content,
		CreatedAt:      row.CreatedAt.Time,
		Sender: model.MessageSender{
			ID:        row.SenderID,
			Username:  textOrEmpty(row.Username),
			AvatarURL: textOrEmpty(row.AvatarUrl),
		},
	}, nil
}
