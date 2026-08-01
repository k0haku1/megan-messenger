package model

import (
	"bytes"
	"megan-messenger/internal/util"
	"time"

	"github.com/google/uuid"
)

type ConversationType string

const (
	ConversationTypeDM    ConversationType = "dm"
	ConversationTypeGroup ConversationType = "group"
)

type Conversation struct {
	ID        uuid.UUID            `json:"id" binding:"required"`
	Type      ConversationType     `json:"type" binding:"required"`
	Title     string               `json:"title,omitempty"`
	Slug      string               `json:"slug,omitempty"`
	Members   []ConversationMember `json:"members,omitempty"`
	CreatedAt time.Time            `json:"createdAt"`
}

func NewGroupConversation(title string) (Conversation, error) {
	slug, err := util.GenerateConversationSlug()
	if err != nil {
		return Conversation{}, err
	}

	return Conversation{
		ID:        uuid.Must(uuid.NewV7()),
		Type:      ConversationTypeGroup,
		Title:     title,
		Slug:      slug,
		CreatedAt: time.Now(),
	}, nil
}

func NewDMConversation() Conversation {
	return Conversation{
		ID:        uuid.Must(uuid.NewV7()),
		Type:      ConversationTypeDM,
		CreatedAt: time.Now(),
	}
}

type ConversationMember struct {
	ID             uuid.UUID `json:"id" binding:"required"`
	UserID         uuid.UUID `json:"userId" binding:"required"`
	ConversationID uuid.UUID `json:"conversationId" binding:"required"`
	JoinedAt       time.Time `json:"joinedAt"`
}

func NewConversationMember(userID, conversationID uuid.UUID) ConversationMember {
	return ConversationMember{
		ID:             uuid.Must(uuid.NewV7()),
		UserID:         userID,
		ConversationID: conversationID,
		JoinedAt:       time.Now(),
	}
}

func DMPairKey(a, b uuid.UUID) (low, high uuid.UUID) {
	if bytes.Compare(a[:], b[:]) < 0 {
		return a, b
	}
	return b, a
}
