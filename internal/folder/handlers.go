package folder

import (
	"errors"
	"megan-messenger/internal/httputil"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	validator *httputil.Validator
	service   *Service
}

func NewHandler(validator *httputil.Validator, service *Service) *Handler {
	return &Handler{validator: validator, service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	folders, err := h.service.List(r.Context(), user.ID)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, map[string]any{"folders": folders})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	var req CreateFolderRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if errs := h.validator.Validate(req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	f, err := h.service.Create(r.Context(), user.ID, req.Name, req.Icon)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidName):
			httputil.BadRequest(w, "Invalid folder name")
		case errors.Is(err, ErrFolderLimitReached):
			httputil.BadRequest(w, "Folder limit reached")
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.Created(w, map[string]any{"folder": f})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		httputil.BadRequest(w, "Invalid folder ID")
		return
	}

	var req UpdateFolderRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if errs := h.validator.Validate(req); errs != nil {
		httputil.ValidationError(w, errs)
		return
	}

	f, err := h.service.Update(r.Context(), user.ID, folderID, req.Name, req.Icon)
	if err != nil {
		switch {
		case errors.Is(err, ErrFolderNotFound):
			httputil.NotFound(w, "Folder not found")
		case errors.Is(err, ErrInvalidName):
			httputil.BadRequest(w, "Invalid folder name")
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.SuccessData(w, map[string]any{"folder": f})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		httputil.BadRequest(w, "Invalid folder ID")
		return
	}

	if err := h.service.Delete(r.Context(), user.ID, folderID); err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			httputil.NotFound(w, "Folder not found")
		} else {
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) Reorder(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)

	var req ReorderRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}

	if err := h.service.Reorder(r.Context(), user.ID, req.FolderIDs); err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) AddItems(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		httputil.BadRequest(w, "Invalid folder ID")
		return
	}

	var req AddItemsRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}

	var firstErr error
	for _, convID := range req.ConversationIDs {
		if err := h.service.AddItem(r.Context(), user.ID, folderID, convID); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if firstErr != nil {
		switch {
		case errors.Is(firstErr, ErrFolderNotFound):
			httputil.NotFound(w, "Folder not found")
		case errors.Is(firstErr, ErrNotMember):
			httputil.Forbidden(w, "Not a member of one or more conversations")
		case errors.Is(firstErr, ErrItemLimitReached):
			httputil.BadRequest(w, "Folder item limit reached")
		default:
			httputil.InternalError(w, r, firstErr)
		}
		return
	}

	httputil.NoContent(w)
}

func (h *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	user := httputil.UserFromRequest(r)
	folderID, err := uuid.Parse(chi.URLParam(r, "folderID"))
	if err != nil {
		httputil.BadRequest(w, "Invalid folder ID")
		return
	}
	convID, err := uuid.Parse(chi.URLParam(r, "conversationID"))
	if err != nil {
		httputil.BadRequest(w, "Invalid conversation ID")
		return
	}

	if err := h.service.RemoveItem(r.Context(), user.ID, folderID, convID); err != nil {
		if errors.Is(err, ErrFolderNotFound) {
			httputil.NotFound(w, "Folder not found")
		} else {
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.NoContent(w)
}
