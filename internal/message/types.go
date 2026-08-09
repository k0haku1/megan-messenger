package message

import (
	"megan-messenger/internal/model"

	"github.com/google/uuid"
)

type GetPagingResponse struct {
	Messages   []model.Message `json:"messages"`
	NextCursor string          `json:"nextCursor"`
}

type SendMessageRequest struct {
	Content   string     `json:"content" validate:"required,min=1,max=5000"`
	ReplyToID *uuid.UUID `json:"replyToId"`
}

type SendMessageResponse struct {
	Message model.Message `json:"message"`
}

type ForwardMessageRequest struct {
	ConversationID uuid.UUID `json:"conversationId" validate:"required"`
}
type ReactionRequest struct {
	Emoji string `json:"emoji" validate:"required,min=1,max=32"`
}
type MessageResponse struct {
	Message model.Message `json:"message"`
}
