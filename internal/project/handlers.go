package project

import (
	"errors"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/model"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	validator *httputil.Validator
	service   *Service
}

func NewHandler(validator *httputil.Validator, service *Service) *Handler {
	return &Handler{
		validator: validator,
		service:   service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	var req CreateProjectRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	project, err := h.service.Create(r.Context(), user.ID, req.Name, req.Description)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.Created(w, ProjectResponse{Project: project})
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	projects, err := h.service.ListMine(r.Context(), user.ID)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	response := make([]ProjectResponse, 0, len(projects))
	for _, item := range projects {
		response = append(response, ProjectResponse{
			Project:     item.Project,
			MemberCount: item.MemberCount,
			DocCount:    item.DocCount,
		})
	}

	httputil.SuccessData(w, ListProjectsResponse{Projects: response})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	project, err := h.service.Get(r.Context(), user.ID, projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, ProjectResponse{Project: project})
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	var req AddMemberRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	if err := h.service.AddMember(r.Context(), user.ID, projectID, req.Username); err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	members, err := h.service.ListMembers(r.Context(), user.ID, projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, MembersResponse{Members: members})
}

func (h *Handler) CreateDoc(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	var req CreateDocRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	doc, err := h.service.CreateDoc(r.Context(), user.ID, projectID, req.Title, req.BodyMD)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.Created(w, DocResponse{Doc: doc})
}

func (h *Handler) ListDocs(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	docs, err := h.service.ListDocs(r.Context(), user.ID, projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, DocsResponse{Docs: docs})
}

func (h *Handler) GetDoc(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}
	docID, err := parseDocID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid doc id")
		return
	}

	doc, err := h.service.GetDoc(r.Context(), user.ID, projectID, docID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, DocResponse{Doc: doc})
}

func (h *Handler) UpdateDoc(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}
	docID, err := parseDocID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid doc id")
		return
	}

	var req UpdateDocRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	doc, err := h.service.UpdateDoc(r.Context(), user.ID, projectID, docID, req.Title, req.BodyMD)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, DocResponse{Doc: doc})
}

func (h *Handler) DeleteDoc(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}
	docID, err := parseDocID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid doc id")
		return
	}

	if err := h.service.DeleteDoc(r.Context(), user.ID, projectID, docID); err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) CreateDecision(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	var req CreateDecisionRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	decision, err := h.service.CreateDecision(r.Context(), user.ID, projectID, req.Summary, req.Context, req.MessageID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.Created(w, DecisionResponse{Decision: decision})
}

func (h *Handler) ListDecisions(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	decisions, err := h.service.ListDecisions(r.Context(), user.ID, projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, DecisionsResponse{Decisions: decisions})
}

func (h *Handler) UpdateDecisionStatus(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}
	decisionID, err := parseDecisionID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid decision id")
		return
	}

	var req UpdateDecisionStatusRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	decision, err := h.service.UpdateDecisionStatus(
		r.Context(),
		user.ID,
		projectID,
		decisionID,
		model.DecisionStatus(req.Status),
	)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, DecisionResponse{Decision: decision})
}

func (h *Handler) LinkConversation(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	var req LinkConversationRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fields := h.validator.Validate(req); len(fields) > 0 {
		httputil.ValidationError(w, fields)
		return
	}

	if err := h.service.LinkConversation(r.Context(), user.ID, projectID, req.ConversationID); err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) UnlinkConversation(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}
	conversationID, err := parseConversationID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid conversation id")
		return
	}

	if err := h.service.UnlinkConversation(r.Context(), user.ID, projectID, conversationID); err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) ListConversations(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	projectID, err := parseProjectID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid project id")
		return
	}

	conversations, err := h.service.ListLinkedConversations(r.Context(), user.ID, projectID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, ConversationsResponse{Conversations: conversations})
}

func (h *Handler) ListConversationProjects(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	conversationID, err := httputil.ParseConversationID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid conversation id")
		return
	}

	projects, err := h.service.ListConversationProjects(r.Context(), user.ID, conversationID)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	httputil.SuccessData(w, ProjectsResponse{Projects: projects})
}

func (h *Handler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrNotFound),
		errors.Is(err, ErrDocNotFound),
		errors.Is(err, ErrDecisionNotFound),
		errors.Is(err, ErrMemberNotFound):
		httputil.NotFound(w, "")
	case errors.Is(err, ErrAccessDenied),
		errors.Is(err, ErrOwnerRequired),
		errors.Is(err, ErrNotConversationMember):
		httputil.Forbidden(w, "")
	case errors.Is(err, ErrNotGroupChat):
		httputil.BadRequest(w, "only group conversations can be linked")
	case errors.Is(err, ErrMessageNotFound):
		httputil.NotFound(w, "message not found")
	case errors.Is(err, ErrMessageNotLinked):
		httputil.BadRequest(w, "message is not from a linked conversation")
	default:
		httputil.InternalError(w, r, err)
	}
}

func parseProjectID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("projectID"))
}

func parseDocID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("docID"))
}

func parseDecisionID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("decisionID"))
}

func parseConversationID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("conversationID"))
}
