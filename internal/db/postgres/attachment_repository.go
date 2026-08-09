package postgres

import (
	"context"
	db "megan-messenger/internal/db/postgres/sqlc"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AttachmentRepository struct {
	queries *db.Queries
}

func NewAttachmentRepository(queries *db.Queries) *AttachmentRepository {
	return &AttachmentRepository{queries: queries}
}

func mapAttachment(row db.MessageAttachment) model.MessageAttachment {
	return model.MessageAttachment{
		ID:             row.ID,
		MessageID:      uuidPtr(row.MessageID),
		ConversationID: row.ConversationID,
		UploaderID:     row.UploaderID,
		ObjectKey:      row.ObjectKey,
		ThumbKey:       textOrEmpty(row.ThumbKey),
		Mime:           row.Mime,
		Kind:           model.AttachmentKind(row.Kind),
		SizeBytes:      row.SizeBytes,
		Width:          int4Ptr(row.Width),
		Height:         int4Ptr(row.Height),
		DurationMs:     int4Ptr(row.DurationMs),
		OriginalName:   row.OriginalName,
		CreatedAt:      row.CreatedAt.Time,
	}
}

func (r *AttachmentRepository) Create(ctx context.Context, att model.MessageAttachment) (model.MessageAttachment, error) {
	created, err := r.queries.CreateAttachment(ctx, db.CreateAttachmentParams{
		ID:             att.ID,
		MessageID:      uuidFromPtr(att.MessageID),
		ConversationID: att.ConversationID,
		UploaderID:     att.UploaderID,
		ObjectKey:      att.ObjectKey,
		ThumbKey:       textFromString(att.ThumbKey),
		Mime:           att.Mime,
		Kind:           string(att.Kind),
		SizeBytes:      att.SizeBytes,
		Width:          int4FromPtr(att.Width),
		Height:         int4FromPtr(att.Height),
		DurationMs:     int4FromPtr(att.DurationMs),
		OriginalName:   att.OriginalName,
		CreatedAt:      timestampFromTime(att.CreatedAt),
	})
	if err != nil {
		return model.MessageAttachment{}, err
	}
	return mapAttachment(created), nil
}

func (r *AttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (model.MessageAttachment, error) {
	row, err := r.queries.GetAttachmentByID(ctx, id)
	if err != nil {
		return model.MessageAttachment{}, err
	}
	return mapAttachment(row), nil
}

func (r *AttachmentRepository) ListByMessageIDs(ctx context.Context, messageIDs []uuid.UUID) ([]model.MessageAttachment, error) {
	if len(messageIDs) == 0 {
		return []model.MessageAttachment{}, nil
	}
	rows, err := r.queries.ListAttachmentsByMessageIDs(ctx, messageIDs)
	if err != nil {
		return nil, err
	}
	out := make([]model.MessageAttachment, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapAttachment(row))
	}
	return out, nil
}

func (r *AttachmentRepository) AttachToMessage(
	ctx context.Context,
	messageID, uploaderID, conversationID uuid.UUID,
	attachmentIDs []uuid.UUID,
) error {
	if len(attachmentIDs) == 0 {
		return nil
	}
	return r.queries.AttachToMessage(ctx, db.AttachToMessageParams{
		MessageID:      pgtype.UUID{Bytes: messageID, Valid: true},
		AttachmentIds:  attachmentIDs,
		UploaderID:     uploaderID,
		ConversationID: conversationID,
	})
}

func (r *AttachmentRepository) CountPending(
	ctx context.Context,
	uploaderID, conversationID uuid.UUID,
	attachmentIDs []uuid.UUID,
) (int64, error) {
	if len(attachmentIDs) == 0 {
		return 0, nil
	}
	return r.queries.CountPendingAttachments(ctx, db.CountPendingAttachmentsParams{
		AttachmentIds:  attachmentIDs,
		UploaderID:     uploaderID,
		ConversationID: conversationID,
	})
}

func (r *AttachmentRepository) ListConversationMedia(
	ctx context.Context,
	conversationID uuid.UUID,
	kind string,
	limit int,
	cursor *pagination.Cursor,
) ([]model.MessageAttachment, error) {
	params := db.ListConversationMediaParams{
		ConversationID: conversationID,
		Kind:           kind,
		Limit:          int32(limit),
	}
	if cursor != nil {
		params.CursorID = cursor.ID
		params.CursorCreatedAt = pgtype.Timestamptz{Time: cursor.CreatedAt, Valid: true}
	}
	rows, err := r.queries.ListConversationMedia(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]model.MessageAttachment, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapAttachment(row))
	}
	return out, nil
}
