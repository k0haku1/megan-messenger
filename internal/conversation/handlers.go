package conversation

import (
	"errors"
	"log/slog"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/repository"
	"megan-messenger/internal/ws"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	validator *httputil.Validator
	service   *Service
	wsService *ws.Service
}

func NewHandler(validator *httputil.Validator, service *Service, wsService *ws.Service) *Handler {
	return &Handler{
		validator: validator,
		service:   service,
		wsService: wsService,
	}
}

// List godoc
//
//	@Summary	List user conversations
//	@Tags		conversation
//	@Security	BearerAuth
//	@Success	200	{object}	ListResponse
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/conversations [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	conversations, err := h.service.ListByUser(r.Context(), user.ID)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, ListResponse{Conversations: conversations})
}

// CreateGroup godoc
//
//	@Summary	Create a group conversation
//	@Tags		conversation
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		input	body		CreateGroupRequest	true	"Group creation params"
//	@Success	201		{object}	CreateGroupResponse
//	@Failure	400		{object}	httputil.ErrorResponse
//	@Failure	401		{object}	httputil.ErrorResponse
//	@Failure	500		{object}	httputil.ErrorResponse
//	@Router		/conversations/group [post]
func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	var params CreateGroupRequest
	if err := httputil.Read(r, &params); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(params); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	created, err := h.service.CreateGroup(r.Context(), user.ID, params.Title, params.MemberIDs)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.Created(w, CreateGroupResponse{ID: created.ID, Slug: created.Slug})
}

// CreateOrGetDM godoc
//
//	@Summary	Create or get a direct message conversation
//	@Tags		conversation
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		input	body		CreateDMRequest	true	"DM params"
//	@Success	200		{object}	model.Conversation
//	@Failure	400		{object}	httputil.ErrorResponse
//	@Failure	401		{object}	httputil.ErrorResponse
//	@Failure	500		{object}	httputil.ErrorResponse
//	@Router		/conversations/dm [post]
func (h *Handler) CreateOrGetDM(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	var params CreateDMRequest
	if err := httputil.Read(r, &params); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(params); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	conv, err := h.service.GetOrCreateDM(r.Context(), user.ID, params.UserID)
	if err != nil {
		if errors.Is(err, ErrCannotDMYourself) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, conv)
}

// JoinBySlug godoc
//
//	@Summary	Join a group conversation by invite slug
//	@Tags		conversation
//	@Param		slug	path	string	true	"Invite slug"
//	@Security	BearerAuth
//	@Success	200	{object}	JoinBySlugResponse
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Failure	404	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/conversations/join/{slug} [post]
func (h *Handler) JoinBySlug(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	slug := r.PathValue("slug")

	conv, err := h.service.JoinBySlug(r.Context(), user.ID, slug)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			httputil.NotFound(w, "")
			return
		}
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, JoinBySlugResponse{ID: conv.ID})
}

// Websocket godoc
//
//	@Summary		Connect to conversation websocket
//	@Tags			conversation
//	@Param			conversationID	path	string	true	"Conversation ID"
//	@Security		WebSocketQueryAuth
//	@Description	Connect to receive real-time messages in a conversation
//	@Schemes		ws
//	@Failure		401	{object}	httputil.ErrorResponse
//	@Failure		404	{object}	httputil.ErrorResponse
//	@Failure		500	{object}	httputil.ErrorResponse
//	@Router			/conversations/{conversationID}/ws [get]
func (h *Handler) Websocket(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	conversationID, err := uuid.Parse(r.PathValue("conversationID"))
	if err != nil {
		httputil.BadRequest(w, "invalid conversation id")
		return
	}

	conv, err := h.service.EnsureMember(r.Context(), user.ID, conversationID)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			httputil.NotFound(w, "")
			return
		}
		httputil.InternalError(w, r, err)
		return
	}

	if err := h.wsService.HandleWebSocket(w, r, conv, user.ID); err != nil {
		slog.Error("websocket error", "err", err)
	}
}
