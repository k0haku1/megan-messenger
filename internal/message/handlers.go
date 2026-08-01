package message

import (
	"errors"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/pagination"
	"net/http"

	"github.com/google/uuid"
)

type Handler struct {
	validate *httputil.Validator
	service  *Service
}

func NewHandler(validator *httputil.Validator, service *Service) *Handler {
	return &Handler{
		validate: validator,
		service:  service,
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

	conversationID, err := uuid.Parse(r.PathValue("conversationID"))
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
		if errors.Is(err, ErrConversationNotFound) {
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
