package conversation

import (
	"errors"
	"log/slog"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	userpkg "megan-messenger/internal/user"
	"megan-messenger/internal/ws"
	"net/http"
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

// SendDMMessage godoc
//
//	@Summary		Send the first message in a direct conversation
//	@Description	Creates the DM conversation if it does not exist yet
//	@Tags			conversation
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			input	body		SendDMMessageRequest	true	"DM message params"
//	@Success		201		{object}	SendDMMessageResponse
//	@Failure		400		{object}	httputil.ErrorResponse
//	@Failure		401		{object}	httputil.ErrorResponse
//	@Failure		403		{object}	httputil.ErrorResponse
//	@Failure		404		{object}	httputil.ErrorResponse
//	@Failure		422		{object}	httputil.ErrorResponse
//	@Failure		500		{object}	httputil.ErrorResponse
//	@Router			/conversations/dm/messages [post]
func (h *Handler) SendDMMessage(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	var params SendDMMessageRequest
	if err := httputil.Read(r, &params); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}

	if err := h.validator.Validate(params); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	target, fieldErrors, ok := parseDMTarget(params.UserID, params.Username)
	if !ok {
		httputil.ValidationError(w, fieldErrors)
		return
	}

	var conv model.Conversation
	var message model.Message
	var err error
	if target.userID != nil {
		conv, message, err = h.service.SendDMMessage(r.Context(), user.ID, target.userID, nil, params.Content)
	} else {
		conv, message, err = h.service.SendDMMessage(r.Context(), user.ID, nil, target.username, params.Content)
	}
	if err != nil {
		switch {
		case errors.Is(err, ErrCannotDMYourself):
			httputil.BadRequest(w, err.Error())
		case errors.Is(err, ErrCannotMessageUser):
			httputil.Forbidden(w, err.Error())
		case errors.Is(err, userpkg.ErrUserNotDiscoverable):
			httputil.NotFound(w, "")
		case errors.Is(err, model.ErrInvalidMessageContent):
			httputil.ValidationError(w, httputil.FieldErrors{"content": "Message is too long"})
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	if err := h.wsService.PublishMessage(r.Context(), conv.ID, message); err != nil {
		slog.Warn("failed to publish dm message", "err", err)
	}

	httputil.Created(w, SendDMMessageResponse{
		Conversation: conv,
		Message:      message,
	})
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

	conversationID, err := httputil.ParseConversationID(r)
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
