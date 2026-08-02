package project

import (
	"megan-messenger/internal/model"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=2000"`
}

type ProjectResponse struct {
	model.Project
	MemberCount int32 `json:"memberCount,omitempty"`
	DocCount    int32 `json:"docCount,omitempty"`
}

type ListProjectsResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

type AddMemberRequest struct {
	Username string `json:"username" validate:"required,min=5,max=32"`
}

type MembersResponse struct {
	Members []model.ProjectMember `json:"members"`
}

type CreateDocRequest struct {
	Title  string `json:"title" validate:"required,min=1,max=200"`
	BodyMD string `json:"bodyMd" validate:"max=100000"`
}

type UpdateDocRequest struct {
	Title  string `json:"title" validate:"required,min=1,max=200"`
	BodyMD string `json:"bodyMd" validate:"max=100000"`
}

type DocResponse struct {
	Doc model.ProjectDoc `json:"doc"`
}

type DocsResponse struct {
	Docs []model.ProjectDocSummary `json:"docs"`
}

type CreateDecisionRequest struct {
	Summary   string     `json:"summary" validate:"required,min=3,max=300"`
	Context   string     `json:"context" validate:"max=5000"`
	MessageID *uuid.UUID `json:"messageId" validate:"omitempty,uuid"`
}

type UpdateDecisionStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=proposed accepted deprecated"`
}

type DecisionResponse struct {
	Decision model.ProjectDecision `json:"decision"`
}

type DecisionsResponse struct {
	Decisions []model.ProjectDecision `json:"decisions"`
}

type LinkConversationRequest struct {
	ConversationID uuid.UUID `json:"conversationId" validate:"required,uuid"`
}

type ConversationsResponse struct {
	Conversations []model.ProjectConversationLink `json:"conversations"`
}

type ProjectsResponse struct {
	Projects []model.Project `json:"projects"`
}
