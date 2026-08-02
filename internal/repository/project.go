package repository

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectDocNotFound   = errors.New("project doc not found")
	ErrProjectDecisionNotFound = errors.New("project decision not found")
	ErrProjectAccessDenied  = errors.New("project access denied")
	ErrProjectDocSlugExists = errors.New("project doc slug exists")
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project model.Project) (model.Project, error)
	GetProject(ctx context.Context, id uuid.UUID) (model.Project, error)
	ListUserProjects(ctx context.Context, userID uuid.UUID) ([]model.ProjectListItem, error)
	IsProjectMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error)
	GetProjectMemberRole(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMemberRole, error)
	AddProjectMember(ctx context.Context, projectID, userID uuid.UUID, role model.ProjectMemberRole, joinedAt time.Time) error
	ListProjectMembers(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error)

	CreateDoc(ctx context.Context, doc model.ProjectDoc) (model.ProjectDoc, error)
	GetDoc(ctx context.Context, projectID, docID uuid.UUID) (model.ProjectDoc, error)
	ListDocs(ctx context.Context, projectID uuid.UUID) ([]model.ProjectDocSummary, error)
	UpdateDoc(ctx context.Context, doc model.ProjectDoc) (model.ProjectDoc, error)
	DeleteDoc(ctx context.Context, projectID, docID uuid.UUID) error
	DocSlugExists(ctx context.Context, projectID uuid.UUID, slug string) (bool, error)

	CreateDecision(ctx context.Context, decision model.ProjectDecision) (model.ProjectDecision, error)
	ListDecisions(ctx context.Context, projectID uuid.UUID) ([]model.ProjectDecision, error)
	GetDecision(ctx context.Context, projectID, decisionID uuid.UUID) (model.ProjectDecision, error)
	UpdateDecisionStatus(ctx context.Context, projectID, decisionID uuid.UUID, status model.DecisionStatus) (model.ProjectDecision, error)

	LinkConversation(ctx context.Context, projectID, conversationID uuid.UUID, linkedAt time.Time) error
	UnlinkConversation(ctx context.Context, projectID, conversationID uuid.UUID) error
	ListLinkedConversations(ctx context.Context, projectID uuid.UUID) ([]model.ProjectConversationLink, error)
	ListLinkedProjects(ctx context.Context, conversationID, userID uuid.UUID) ([]model.Project, error)
	IsConversationLinked(ctx context.Context, projectID, conversationID uuid.UUID) (bool, error)
}
