package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMemberRole string

const (
	ProjectRoleOwner  ProjectMemberRole = "owner"
	ProjectRoleMember ProjectMemberRole = "member"
)

type DecisionStatus string

const (
	DecisionProposed   DecisionStatus = "proposed"
	DecisionAccepted   DecisionStatus = "accepted"
	DecisionDeprecated DecisionStatus = "deprecated"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ProjectListItem struct {
	Project
	MemberCount int32 `json:"memberCount"`
	DocCount    int32 `json:"docCount"`
}

type ProjectMember struct {
	ProjectID uuid.UUID         `json:"projectId"`
	UserID    uuid.UUID         `json:"userId"`
	Role      ProjectMemberRole `json:"role"`
	JoinedAt  time.Time         `json:"joinedAt"`
	Username  string            `json:"username"`
	AvatarURL string            `json:"avatarUrl,omitempty"`
}

type ProjectDoc struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	BodyMD    string    `json:"bodyMd"`
	CreatedBy uuid.UUID `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProjectDocSummary struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	CreatedBy uuid.UUID `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProjectDecision struct {
	ID        uuid.UUID      `json:"id"`
	ProjectID uuid.UUID      `json:"projectId"`
	Summary   string         `json:"summary"`
	Context   string         `json:"context"`
	Status    DecisionStatus `json:"status"`
	MessageID *uuid.UUID     `json:"messageId,omitempty"`
	CreatedBy uuid.UUID      `json:"createdBy"`
	CreatedAt time.Time      `json:"createdAt"`
}

type ProjectConversationLink struct {
	ProjectID      uuid.UUID        `json:"projectId"`
	ConversationID uuid.UUID        `json:"conversationId"`
	LinkedAt       time.Time        `json:"linkedAt"`
	Type           ConversationType `json:"type"`
	Title          string           `json:"title"`
	Slug           string           `json:"slug,omitempty"`
}
