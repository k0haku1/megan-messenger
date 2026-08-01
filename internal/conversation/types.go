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

type CreateDMRequest struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
}

type ListResponse struct {
	Conversations []model.Conversation `json:"conversations"`
}

type JoinBySlugResponse struct {
	ID uuid.UUID `json:"id"`
}
