package conversation

import (
	"megan-messenger/internal/model"

	"github.com/google/uuid"
)

type CreateGroupRequest struct {
	Title     string      `json:"title" validate:"required,min=3,max=100"`
	MemberIDs []uuid.UUID `json:"memberIds"`
}

type CreateGroupResponse struct {
	ID   uuid.UUID `json:"id"`
	Slug string    `json:"slug"`
}

type SendDMMessageRequest struct {
	UserID   *uuid.UUID `json:"userId" validate:"omitempty,uuid"`
	Username *string    `json:"username" validate:"omitempty,min=5,max=32"`
	Content  string     `json:"content" validate:"required,min=1,max=5000"`
}

type SendDMMessageResponse struct {
	Conversation model.Conversation `json:"conversation"`
	Message      model.Message      `json:"message"`
}

type ListResponse struct {
	Conversations []model.Conversation `json:"conversations"`
}

type JoinBySlugResponse struct {
	ID uuid.UUID `json:"id"`
}
