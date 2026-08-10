package ws

import (
	"megan-messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

const (
	EventTypeMessage      = "message"
	EventTypeRead         = "read"
	EventTypeConversation = "conversation"
)

type Event struct {
	Type         string                   `json:"type"`
	Message      *model.Message           `json:"message,omitempty"`
	Read         *ReadEvent               `json:"read,omitempty"`
	Conversation *ConversationActivityEvent `json:"conversation,omitempty"`
}

type ReadEvent struct {
	ConversationID uuid.UUID `json:"conversationId"`
	UserID         uuid.UUID `json:"userId"`
	LastReadAt     time.Time `json:"lastReadAt"`
}

// ConversationActivityEvent is a lightweight inbox update for the chat list.
type ConversationActivityEvent struct {
	ConversationID uuid.UUID `json:"conversationId"`
	MessageID      uuid.UUID `json:"messageId"`
	Content        string    `json:"content"`
	SenderID       uuid.UUID `json:"senderId"`
	SenderUsername string    `json:"senderUsername"`
	AttachmentKind string    `json:"attachmentKind,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

func MessageEvent(message model.Message) Event {
	msg := message
	return Event{Type: EventTypeMessage, Message: &msg}
}

func ReadReceiptEvent(conversationID, userID uuid.UUID, lastReadAt time.Time) Event {
	return Event{
		Type: EventTypeRead,
		Read: &ReadEvent{
			ConversationID: conversationID,
			UserID:         userID,
			LastReadAt:     lastReadAt,
		},
	}
}

func ConversationActivityFromMessage(message model.Message) Event {
	kind := ""
	if len(message.Attachments) > 0 {
		kind = string(message.Attachments[0].Kind)
	}
	return Event{
		Type: EventTypeConversation,
		Conversation: &ConversationActivityEvent{
			ConversationID: message.ConversationID,
			MessageID:      message.ID,
			Content:        message.Content,
			SenderID:       message.Sender.ID,
			SenderUsername: message.Sender.Username,
			AttachmentKind: kind,
			CreatedAt:      message.CreatedAt,
		},
	}
}

func UserChannel(userID uuid.UUID) string {
	return "user:" + userID.String()
}
