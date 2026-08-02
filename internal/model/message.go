package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidMessageContent = errors.New("invalid message content")

type MessageSender struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatarUrl"`
}

type Message struct {
	ID             uuid.UUID     `json:"id" binding:"required"`
	ConversationID uuid.UUID     `json:"conversationId" binding:"required"`
	Content        string        `json:"content" binding:"required"`
	Sender         MessageSender `json:"sender" binding:"required"`
	CreatedAt      time.Time     `json:"createdAt" binding:"required"`
}

func NewMessage(conversationID uuid.UUID, content string, sender User) (Message, error) {
	if len(content) > 5000 {
		return Message{}, ErrInvalidMessageContent
	}
	return Message{
		ID:             uuid.Must(uuid.NewV7()),
		ConversationID: conversationID,
		Content:        content,
		Sender: MessageSender{
			ID:        sender.ID,
			Username:  sender.Username, // empty until onboarding completes
			AvatarURL: sender.AvatarURL,
		},
		CreatedAt: time.Now(),
	}, nil
}
