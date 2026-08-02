package postgres

import (
	"context"
	"errors"
	db "megan-messenger/internal/db/postgres/sqlc"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectRepository struct {
	queries db.Querier
}

func NewProjectRepository(queries db.Querier) repository.ProjectRepository {
	return &ProjectRepository{queries: queries}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project model.Project) (model.Project, error) {
	row, err := r.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   timestampFromTime(project.CreatedAt),
	})
	if err != nil {
		return model.Project{}, err
	}
	return mapProject(row), nil
}

func (r *ProjectRepository) GetProject(ctx context.Context, id uuid.UUID) (model.Project, error) {
	row, err := r.queries.GetProject(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Project{}, repository.ErrProjectNotFound
		}
		return model.Project{}, err
	}
	return mapProject(row), nil
}

func (r *ProjectRepository) ListUserProjects(ctx context.Context, userID uuid.UUID) ([]model.ProjectListItem, error) {
	rows, err := r.queries.ListUserProjects(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]model.ProjectListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.ProjectListItem{
			Project: model.Project{
				ID:          row.ID,
				Name:        row.Name,
				Description: row.Description,
				CreatedBy:   row.CreatedBy,
				CreatedAt:   row.CreatedAt.Time,
			},
			MemberCount: row.MemberCount,
			DocCount:    row.DocCount,
		})
	}
	return items, nil
}

func (r *ProjectRepository) IsProjectMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	return r.queries.IsProjectMember(ctx, db.IsProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
	})
}

func (r *ProjectRepository) GetProjectMemberRole(ctx context.Context, projectID, userID uuid.UUID) (model.ProjectMemberRole, error) {
	role, err := r.queries.GetProjectMemberRole(ctx, db.GetProjectMemberRoleParams{
		ProjectID: projectID,
		UserID:    userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrProjectAccessDenied
		}
		return "", err
	}
	return model.ProjectMemberRole(role), nil
}

func (r *ProjectRepository) AddProjectMember(
	ctx context.Context,
	projectID, userID uuid.UUID,
	role model.ProjectMemberRole,
	joinedAt time.Time,
) error {
	return r.queries.AddProjectMember(ctx, db.AddProjectMemberParams{
		ProjectID: projectID,
		UserID:    userID,
		Role:      db.ProjectMemberRole(role),
		JoinedAt:  timestampFromTime(joinedAt),
	})
}

func (r *ProjectRepository) ListProjectMembers(ctx context.Context, projectID uuid.UUID) ([]model.ProjectMember, error) {
	rows, err := r.queries.ListProjectMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}
	members := make([]model.ProjectMember, 0, len(rows))
	for _, row := range rows {
		members = append(members, model.ProjectMember{
			ProjectID: row.ProjectID,
			UserID:    row.UserID,
			Role:      model.ProjectMemberRole(row.Role),
			JoinedAt:  row.JoinedAt.Time,
			Username:  textOrEmpty(row.Username),
			AvatarURL: textOrEmpty(row.AvatarUrl),
		})
	}
	return members, nil
}

func (r *ProjectRepository) CreateDoc(ctx context.Context, doc model.ProjectDoc) (model.ProjectDoc, error) {
	row, err := r.queries.CreateProjectDoc(ctx, db.CreateProjectDocParams{
		ID:        doc.ID,
		ProjectID: doc.ProjectID,
		Slug:      doc.Slug,
		Title:     doc.Title,
		BodyMd:    doc.BodyMD,
		CreatedBy: doc.CreatedBy,
		UpdatedAt: timestampFromTime(doc.UpdatedAt),
	})
	if err != nil {
		return model.ProjectDoc{}, err
	}
	return mapProjectDoc(row), nil
}

func (r *ProjectRepository) GetDoc(ctx context.Context, projectID, docID uuid.UUID) (model.ProjectDoc, error) {
	row, err := r.queries.GetProjectDoc(ctx, db.GetProjectDocParams{
		ID:        docID,
		ProjectID: projectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectDoc{}, repository.ErrProjectDocNotFound
		}
		return model.ProjectDoc{}, err
	}
	return mapProjectDoc(row), nil
}

func (r *ProjectRepository) ListDocs(ctx context.Context, projectID uuid.UUID) ([]model.ProjectDocSummary, error) {
	rows, err := r.queries.ListProjectDocs(ctx, projectID)
	if err != nil {
		return nil, err
	}
	docs := make([]model.ProjectDocSummary, 0, len(rows))
	for _, row := range rows {
		docs = append(docs, model.ProjectDocSummary{
			ID:        row.ID,
			ProjectID: row.ProjectID,
			Slug:      row.Slug,
			Title:     row.Title,
			CreatedBy: row.CreatedBy,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return docs, nil
}

func (r *ProjectRepository) UpdateDoc(ctx context.Context, doc model.ProjectDoc) (model.ProjectDoc, error) {
	row, err := r.queries.UpdateProjectDoc(ctx, db.UpdateProjectDocParams{
		ID:        doc.ID,
		ProjectID: doc.ProjectID,
		Title:     doc.Title,
		BodyMd:    doc.BodyMD,
		UpdatedAt: timestampFromTime(doc.UpdatedAt),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectDoc{}, repository.ErrProjectDocNotFound
		}
		return model.ProjectDoc{}, err
	}
	return mapProjectDoc(row), nil
}

func (r *ProjectRepository) DeleteDoc(ctx context.Context, projectID, docID uuid.UUID) error {
	return r.queries.DeleteProjectDoc(ctx, db.DeleteProjectDocParams{
		ID:        docID,
		ProjectID: projectID,
	})
}

func (r *ProjectRepository) DocSlugExists(ctx context.Context, projectID uuid.UUID, slug string) (bool, error) {
	return r.queries.ProjectDocSlugExists(ctx, db.ProjectDocSlugExistsParams{
		ProjectID: projectID,
		Slug:      slug,
	})
}

func (r *ProjectRepository) CreateDecision(ctx context.Context, decision model.ProjectDecision) (model.ProjectDecision, error) {
	row, err := r.queries.CreateProjectDecision(ctx, db.CreateProjectDecisionParams{
		ID:        decision.ID,
		ProjectID: decision.ProjectID,
		Summary:   decision.Summary,
		Context:   decision.Context,
		Status:    db.DecisionStatus(decision.Status),
		MessageID: uuidToPg(decision.MessageID),
		CreatedBy: decision.CreatedBy,
		CreatedAt: timestampFromTime(decision.CreatedAt),
	})
	if err != nil {
		return model.ProjectDecision{}, err
	}
	return mapProjectDecision(row), nil
}

func (r *ProjectRepository) ListDecisions(ctx context.Context, projectID uuid.UUID) ([]model.ProjectDecision, error) {
	rows, err := r.queries.ListProjectDecisions(ctx, projectID)
	if err != nil {
		return nil, err
	}
	decisions := make([]model.ProjectDecision, 0, len(rows))
	for _, row := range rows {
		decisions = append(decisions, mapProjectDecision(row))
	}
	return decisions, nil
}

func (r *ProjectRepository) GetDecision(ctx context.Context, projectID, decisionID uuid.UUID) (model.ProjectDecision, error) {
	row, err := r.queries.GetProjectDecision(ctx, db.GetProjectDecisionParams{
		ID:        decisionID,
		ProjectID: projectID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectDecision{}, repository.ErrProjectDecisionNotFound
		}
		return model.ProjectDecision{}, err
	}
	return mapProjectDecision(row), nil
}

func (r *ProjectRepository) UpdateDecisionStatus(
	ctx context.Context,
	projectID, decisionID uuid.UUID,
	status model.DecisionStatus,
) (model.ProjectDecision, error) {
	row, err := r.queries.UpdateProjectDecisionStatus(ctx, db.UpdateProjectDecisionStatusParams{
		ID:        decisionID,
		ProjectID: projectID,
		Status:    db.DecisionStatus(status),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ProjectDecision{}, repository.ErrProjectDecisionNotFound
		}
		return model.ProjectDecision{}, err
	}
	return mapProjectDecision(row), nil
}

func (r *ProjectRepository) LinkConversation(
	ctx context.Context,
	projectID, conversationID uuid.UUID,
	linkedAt time.Time,
) error {
	return r.queries.LinkProjectConversation(ctx, db.LinkProjectConversationParams{
		ProjectID:      projectID,
		ConversationID: conversationID,
		LinkedAt:       timestampFromTime(linkedAt),
	})
}

func (r *ProjectRepository) UnlinkConversation(ctx context.Context, projectID, conversationID uuid.UUID) error {
	return r.queries.UnlinkProjectConversation(ctx, db.UnlinkProjectConversationParams{
		ProjectID:      projectID,
		ConversationID: conversationID,
	})
}

func (r *ProjectRepository) ListLinkedConversations(ctx context.Context, projectID uuid.UUID) ([]model.ProjectConversationLink, error) {
	rows, err := r.queries.ListProjectConversations(ctx, projectID)
	if err != nil {
		return nil, err
	}
	links := make([]model.ProjectConversationLink, 0, len(rows))
	for _, row := range rows {
		links = append(links, model.ProjectConversationLink{
			ProjectID:      row.ProjectID,
			ConversationID: row.ConversationID,
			LinkedAt:       row.LinkedAt.Time,
			Type:           model.ConversationType(row.Type),
			Title:          textOrEmpty(row.Title),
			Slug:           textOrEmpty(row.Slug),
		})
	}
	return links, nil
}

func (r *ProjectRepository) ListLinkedProjects(ctx context.Context, conversationID, userID uuid.UUID) ([]model.Project, error) {
	rows, err := r.queries.ListConversationProjects(ctx, db.ListConversationProjectsParams{
		ConversationID: conversationID,
		UserID:         userID,
	})
	if err != nil {
		return nil, err
	}
	projects := make([]model.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, mapProject(row))
	}
	return projects, nil
}

func (r *ProjectRepository) IsConversationLinked(ctx context.Context, projectID, conversationID uuid.UUID) (bool, error) {
	return r.queries.IsProjectConversationLinked(ctx, db.IsProjectConversationLinkedParams{
		ProjectID:      projectID,
		ConversationID: conversationID,
	})
}

func mapProject(row db.Project) model.Project {
	return model.Project{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		CreatedBy:   row.CreatedBy,
		CreatedAt:   row.CreatedAt.Time,
	}
}

func mapProjectDoc(row db.ProjectDoc) model.ProjectDoc {
	return model.ProjectDoc{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Slug:      row.Slug,
		Title:     row.Title,
		BodyMD:    row.BodyMd,
		CreatedBy: row.CreatedBy,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func mapProjectDecision(row db.ProjectDecision) model.ProjectDecision {
	var messageID *uuid.UUID
	if row.MessageID.Valid {
		id := uuid.UUID(row.MessageID.Bytes)
		messageID = &id
	}
	return model.ProjectDecision{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Summary:   row.Summary,
		Context:   row.Context,
		Status:    model.DecisionStatus(row.Status),
		MessageID: messageID,
		CreatedBy: row.CreatedBy,
		CreatedAt: row.CreatedAt.Time,
	}
}

func uuidToPg(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}
