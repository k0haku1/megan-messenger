package message

import (
	"megan-messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

type GetPagingResponse struct {
	Messages     []model.Message `json:"messages"`
	NextCursor   string          `json:"nextCursor"`
	OthersReadAt *time.Time      `json:"othersReadAt,omitempty"`
}

type SendMessageRequest struct {
	Content       string      `json:"content" validate:"omitempty,max=5000"`
	ReplyToID     *uuid.UUID  `json:"replyToId"`
	AttachmentIDs []uuid.UUID `json:"attachmentIds"`
}

type UploadAttachmentResponse struct {
	Attachment model.MessageAttachment `json:"attachment"`
}

type MediaListResponse struct {
	Items      []model.MessageAttachment `json:"items"`
	NextCursor string                    `json:"nextCursor"`
}

type SendMessageResponse struct {
	Message model.Message `json:"message"`
}

type ForwardMessageRequest struct {
	ConversationID uuid.UUID `json:"conversationId" validate:"required"`
}

type ForwardMessagesRequest struct {
	ConversationID uuid.UUID   `json:"conversationId" validate:"required"`
	MessageIDs     []uuid.UUID `json:"messageIds" validate:"required,min=1,max=50,dive,required"`
}

type ForwardMessagesResponse struct {
	Messages []model.Message `json:"messages"`
}

type ReactionRequest struct {
	Emoji string `json:"emoji" validate:"required,min=1,max=32"`
}
type MessageResponse struct {
	Message model.Message `json:"message"`
}
