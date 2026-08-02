package project

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	projects      repository.ProjectRepository
	users         repository.UserRepository
	conversations repository.ConversationRepository
	messages      repository.MessageRepository
}

func NewService(
	projects repository.ProjectRepository,
	users repository.UserRepository,
	conversations repository.ConversationRepository,
	messages repository.MessageRepository,
) *Service {
	return &Service{
		projects:      projects,
		users:         users,
		conversations: conversations,
		messages:      messages,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, name, description string) (model.Project, error) {
	now := time.Now().UTC()
	project := model.Project{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		CreatedBy:   userID,
		CreatedAt:   now,
	}

	created, err := s.projects.CreateProject(ctx, project)
	if err != nil {
		return model.Project{}, err
	}

	if err := s.projects.AddProjectMember(ctx, created.ID, userID, model.ProjectRoleOwner, now); err != nil {
		return model.Project{}, err
	}

	return created, nil
}

func (s *Service) ListMine(ctx context.Context, userID uuid.UUID) ([]model.ProjectListItem, error) {
	return s.projects.ListUserProjects(ctx, userID)
}

func (s *Service) Get(ctx context.Context, userID, projectID uuid.UUID) (model.Project, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return model.Project{}, err
	}
	return s.projects.GetProject(ctx, projectID)
}

func (s *Service) AddMember(ctx context.Context, actorID, projectID uuid.UUID, username string) error {
	if err := s.requireOwner(ctx, projectID, actorID); err != nil {
		return err
	}

	user, err := s.users.GetByUsername(ctx, strings.TrimPrefix(strings.TrimSpace(username), "@"))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrMemberNotFound
		}
		return err
	}

	if !user.OnboardingComplete() {
		return ErrMemberNotFound
	}

	return s.projects.AddProjectMember(ctx, projectID, user.ID, model.ProjectRoleMember, time.Now().UTC())
}

func (s *Service) ListMembers(ctx context.Context, userID, projectID uuid.UUID) ([]model.ProjectMember, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.projects.ListProjectMembers(ctx, projectID)
}

func (s *Service) CreateDoc(ctx context.Context, userID, projectID uuid.UUID, title, bodyMD string) (model.ProjectDoc, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return model.ProjectDoc{}, err
	}

	now := time.Now().UTC()
	slug, err := s.uniqueSlug(ctx, projectID, SlugFromTitle(title))
	if err != nil {
		return model.ProjectDoc{}, err
	}

	return s.projects.CreateDoc(ctx, model.ProjectDoc{
		ID:        uuid.New(),
		ProjectID: projectID,
		Slug:      slug,
		Title:     strings.TrimSpace(title),
		BodyMD:    bodyMD,
		CreatedBy: userID,
		UpdatedAt: now,
	})
}

func (s *Service) ListDocs(ctx context.Context, userID, projectID uuid.UUID) ([]model.ProjectDocSummary, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.projects.ListDocs(ctx, projectID)
}

func (s *Service) GetDoc(ctx context.Context, userID, projectID, docID uuid.UUID) (model.ProjectDoc, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return model.ProjectDoc{}, err
	}
	return s.projects.GetDoc(ctx, projectID, docID)
}

func (s *Service) UpdateDoc(
	ctx context.Context,
	userID, projectID, docID uuid.UUID,
	title, bodyMD string,
) (model.ProjectDoc, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return model.ProjectDoc{}, err
	}

	existing, err := s.projects.GetDoc(ctx, projectID, docID)
	if err != nil {
		return model.ProjectDoc{}, err
	}

	existing.Title = strings.TrimSpace(title)
	existing.BodyMD = bodyMD
	existing.UpdatedAt = time.Now().UTC()

	return s.projects.UpdateDoc(ctx, existing)
}

func (s *Service) DeleteDoc(ctx context.Context, userID, projectID, docID uuid.UUID) error {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return err
	}
	return s.projects.DeleteDoc(ctx, projectID, docID)
}

func (s *Service) CreateDecision(
	ctx context.Context,
	userID, projectID uuid.UUID,
	summary, contextText string,
	messageID *uuid.UUID,
) (model.ProjectDecision, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return model.ProjectDecision{}, err
	}

	if messageID != nil {
		if err := s.validateDecisionMessage(ctx, userID, projectID, *messageID); err != nil {
			return model.ProjectDecision{}, err
		}
	}

	return s.projects.CreateDecision(ctx, model.ProjectDecision{
		ID:        uuid.New(),
		ProjectID: projectID,
		Summary:   strings.TrimSpace(summary),
		Context:   strings.TrimSpace(contextText),
		Status:    model.DecisionProposed,
		MessageID: messageID,
		CreatedBy: userID,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *Service) ListDecisions(ctx context.Context, userID, projectID uuid.UUID) ([]model.ProjectDecision, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.projects.ListDecisions(ctx, projectID)
}

func (s *Service) UpdateDecisionStatus(
	ctx context.Context,
	userID, projectID, decisionID uuid.UUID,
	status model.DecisionStatus,
) (model.ProjectDecision, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return model.ProjectDecision{}, err
	}
	return s.projects.UpdateDecisionStatus(ctx, projectID, decisionID, status)
}

func (s *Service) LinkConversation(
	ctx context.Context,
	userID, projectID, conversationID uuid.UUID,
) error {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return err
	}

	conv, err := s.conversations.GetByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			return ErrNotConversationMember
		}
		return err
	}

	if conv.Type != model.ConversationTypeGroup {
		return ErrNotGroupChat
	}

	isMember, err := s.conversations.IsMember(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotConversationMember
	}

	return s.projects.LinkConversation(ctx, projectID, conversationID, time.Now().UTC())
}

func (s *Service) UnlinkConversation(
	ctx context.Context,
	userID, projectID, conversationID uuid.UUID,
) error {
	if err := s.requireOwner(ctx, projectID, userID); err != nil {
		return err
	}
	return s.projects.UnlinkConversation(ctx, projectID, conversationID)
}

func (s *Service) ListLinkedConversations(
	ctx context.Context,
	userID, projectID uuid.UUID,
) ([]model.ProjectConversationLink, error) {
	if err := s.requireMember(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.projects.ListLinkedConversations(ctx, projectID)
}

func (s *Service) ListConversationProjects(
	ctx context.Context,
	userID, conversationID uuid.UUID,
) ([]model.Project, error) {
	isMember, err := s.conversations.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrNotConversationMember
	}
	return s.projects.ListLinkedProjects(ctx, conversationID, userID)
}

func (s *Service) validateDecisionMessage(
	ctx context.Context,
	userID, projectID, messageID uuid.UUID,
) error {
	msg, err := s.messages.GetByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			return ErrMessageNotFound
		}
		return err
	}

	linked, err := s.projects.IsConversationLinked(ctx, projectID, msg.ConversationID)
	if err != nil {
		return err
	}
	if !linked {
		return ErrMessageNotLinked
	}

	isMember, err := s.conversations.IsMember(ctx, msg.ConversationID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotConversationMember
	}

	return nil
}

func (s *Service) requireMember(ctx context.Context, projectID, userID uuid.UUID) error {
	ok, err := s.projects.IsProjectMember(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if !ok {
		if _, getErr := s.projects.GetProject(ctx, projectID); getErr != nil {
			if errors.Is(getErr, repository.ErrProjectNotFound) {
				return ErrNotFound
			}
			return getErr
		}
		return ErrAccessDenied
	}
	return nil
}

func (s *Service) requireOwner(ctx context.Context, projectID, userID uuid.UUID) error {
	role, err := s.projects.GetProjectMemberRole(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, repository.ErrProjectAccessDenied) {
			if _, getErr := s.projects.GetProject(ctx, projectID); getErr != nil {
				if errors.Is(getErr, repository.ErrProjectNotFound) {
					return ErrNotFound
				}
				return getErr
			}
			return ErrAccessDenied
		}
		return err
	}
	if role != model.ProjectRoleOwner {
		return ErrOwnerRequired
	}
	return nil
}

func (s *Service) uniqueSlug(ctx context.Context, projectID uuid.UUID, base string) (string, error) {
	candidate := base
	for i := 0; i < 50; i++ {
		exists, err := s.projects.DocSlugExists(ctx, projectID, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = base + "-" + uuid.New().String()[:8]
	}
	return "", repository.ErrProjectDocSlugExists
}
