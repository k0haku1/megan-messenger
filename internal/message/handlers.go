package message

import (
	"errors"
	"log/slog"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"
	"megan-messenger/internal/repository"
	"megan-messenger/internal/ws"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	validate  *httputil.Validator
	service   *Service
	wsService *ws.Service
}

func NewHandler(validator *httputil.Validator, service *Service, wsService *ws.Service) *Handler {
	return &Handler{
		validate:  validator,
		service:   service,
		wsService: wsService,
	}
}

// ListMessages godoc
//
//	@Summary	List messages in a conversation
//	@Tags		message
//	@Produce	json
//	@Security	BearerAuth
//	@Param		conversationID	path		string	true	"Conversation ID"
//	@Param		limit			query		int		false	"Limit"
//	@Param		cursor			query		string	false	"Cursor"
//	@Success	200				{object}	GetPagingResponse
//	@Failure	400				{object}	httputil.ErrorResponse
//	@Failure	404				{object}	httputil.ErrorResponse
//	@Failure	500				{object}	httputil.ErrorResponse
//	@Router		/conversations/{conversationID}/messages [get]
func (h *Handler) ListMessages(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	conversationID, err := httputil.ParseConversationID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid conversation id")
		return
	}

	limit := normalizeLimit(r.URL.Query().Get("limit"), 100, 32)

	var cursor *pagination.Cursor
	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		c, err := h.service.ParseCursor(cursorStr)
		if err != nil {
			httputil.BadRequest(w, "invalid cursor")
			return
		}
		cursor = &c
	}

	messages, err := h.service.ListMessages(r.Context(), user.ID, conversationID, limit, cursor)
	if err != nil {
		if errors.Is(err, repository.ErrConversationNotFound) {
			httputil.NotFound(w, "")
			return
		}
		httputil.InternalError(w, r, err)
		return
	}

	var nextCursor string
	if len(messages) == limit {
		nextCursor = h.service.GenerateCursor(messages[len(messages)-1])
	}

	httputil.SuccessData(w, GetPagingResponse{
		Messages:   messages,
		NextCursor: nextCursor,
	})
}

// SendMessage godoc
//
//	@Summary	Send a message in an existing conversation
//	@Tags		message
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		conversationID	path		string				true	"Conversation ID"
//	@Param		input			body		SendMessageRequest	true	"Message params"
//	@Success	201				{object}	SendMessageResponse
//	@Failure	400				{object}	httputil.ErrorResponse
//	@Failure	401				{object}	httputil.ErrorResponse
//	@Failure	404				{object}	httputil.ErrorResponse
//	@Failure	422				{object}	httputil.ErrorResponse
//	@Failure	500				{object}	httputil.ErrorResponse
//	@Router		/conversations/{conversationID}/messages [post]
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	conversationID, err := httputil.ParseConversationID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid conversation id")
		return
	}

	var params SendMessageRequest
	if err := httputil.Read(r, &params); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}

	if err := h.validate.Validate(params); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	message, err := h.service.SendMessage(r.Context(), user.ID, conversationID, params.Content, params.ReplyToID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrConversationNotFound):
			httputil.NotFound(w, "")
		case errors.Is(err, model.ErrInvalidMessageContent):
			httputil.ValidationError(w, httputil.FieldErrors{"content": "Message is too long"})
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	if err := h.wsService.PublishMessage(r.Context(), conversationID, message); err != nil {
		slog.Warn("failed to publish message", "err", err)
	}

	httputil.Created(w, SendMessageResponse{Message: message})
}

func parseMessageID(r *http.Request) (uuid.UUID, error) { return uuid.Parse(r.PathValue("messageID")) }

func (h *Handler) HideForMe(w http.ResponseWriter, r *http.Request) {
	messageID, err := parseMessageID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid message id")
		return
	}
	err = h.service.HideForUser(r.Context(), httputil.UserFromRequest(r).ID, messageID)
	if respondMessageActionError(w, r, err) {
		return
	}
	httputil.NoContent(w)
}

func (h *Handler) DeleteForEveryone(w http.ResponseWriter, r *http.Request) {
	messageID, err := parseMessageID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid message id")
		return
	}
	message, err := h.service.DeleteForEveryone(r.Context(), httputil.UserFromRequest(r).ID, messageID)
	if respondMessageActionError(w, r, err) {
		return
	}
	if err := h.wsService.PublishMessage(r.Context(), message.ConversationID, message); err != nil {
		slog.Warn("failed to publish message update", "err", err)
	}
	httputil.SuccessData(w, MessageResponse{Message: message})
}

func (h *Handler) ForwardMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := parseMessageID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid message id")
		return
	}
	var params ForwardMessageRequest
	if err := httputil.Read(r, &params); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if err := h.validate.Validate(params); err != nil {
		httputil.ValidationError(w, err)
		return
	}
	message, err := h.service.ForwardMessage(r.Context(), httputil.UserFromRequest(r).ID, messageID, params.ConversationID)
	if respondMessageActionError(w, r, err) {
		return
	}
	if err := h.wsService.PublishMessage(r.Context(), message.ConversationID, message); err != nil {
		slog.Warn("failed to publish forwarded message", "err", err)
	}
	httputil.Created(w, MessageResponse{Message: message})
}

func (h *Handler) AddReaction(w http.ResponseWriter, r *http.Request) { h.changeReaction(w, r, true) }
func (h *Handler) RemoveReaction(w http.ResponseWriter, r *http.Request) {
	h.changeReaction(w, r, false)
}

func (h *Handler) changeReaction(w http.ResponseWriter, r *http.Request, add bool) {
	messageID, err := parseMessageID(r)
	if err != nil {
		httputil.BadRequest(w, "invalid message id")
		return
	}
	var params ReactionRequest
	if err := httputil.Read(r, &params); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if err := h.validate.Validate(params); err != nil {
		httputil.ValidationError(w, err)
		return
	}
	message, err := h.service.SetReaction(r.Context(), httputil.UserFromRequest(r).ID, messageID, params.Emoji, add)
	if respondMessageActionError(w, r, err) {
		return
	}
	if err := h.wsService.PublishMessage(r.Context(), message.ConversationID, message); err != nil {
		slog.Warn("failed to publish message reaction", "err", err)
	}
	httputil.SuccessData(w, MessageResponse{Message: message})
}

func respondMessageActionError(w http.ResponseWriter, r *http.Request, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, repository.ErrMessageNotFound), errors.Is(err, repository.ErrConversationNotFound):
		httputil.NotFound(w, "")
	case errors.Is(err, ErrNotMessageSender):
		httputil.Forbidden(w, "Only the sender can delete this message")
	case errors.Is(err, model.ErrInvalidMessageContent):
		httputil.ValidationError(w, httputil.FieldErrors{"emoji": "Invalid emoji"})
	default:
		httputil.InternalError(w, r, err)
	}
	return true
}
