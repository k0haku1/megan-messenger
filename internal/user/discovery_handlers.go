package user

import (
	"fmt"
	"megan-messenger/internal/httputil"
	"megan-messenger/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

const searchRateLimit = 30

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	userCtx := httputil.UserFromRequest(r)
	query := r.URL.Query().Get("q")

	allowed, err := h.rateLimiter.Allow(
		r.Context(),
		fmt.Sprintf("user-search:%s", userCtx.ID),
		searchRateLimit,
		time.Minute,
	)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}
	if !allowed {
		httputil.Error(w, http.StatusTooManyRequests, "rate_limit_exceeded", "Too many search requests")
		return
	}

	limit := int32(20)
	if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
		parsed, parseErr := strconv.Atoi(rawLimit)
		if parseErr == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}

	users, err := h.service.SearchUsers(r.Context(), userCtx.ID, query, limit)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, SearchUsersResponse{Users: users})
}

func (h *Handler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	userCtx := httputil.UserFromRequest(r)
	username := chi.URLParam(r, "username")

	profile, err := h.service.GetPublicProfile(r.Context(), userCtx.ID, username)
	if err != nil {
		if err == ErrUserNotDiscoverable {
			httputil.NotFound(w, "")
			return
		}
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, profile)
}

func (h *Handler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	userCtx := httputil.UserFromRequest(r)

	var req UpdateUsernameRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if err := h.validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	user, err := h.service.ChangeUsername(r.Context(), userCtx.ID, req.Username)
	if err != nil {
		switch {
		case err == ErrInvalidUsername:
			httputil.ValidationError(w, httputil.FieldErrors{"username": err.Error()})
		case err == ErrUsernameTaken:
			httputil.ValidationError(w, httputil.FieldErrors{"username": err.Error()})
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.SuccessData(w, UpdateUsernameResponse{Username: user.Username})
}

func (h *Handler) UpdatePrivacy(w http.ResponseWriter, r *http.Request) {
	userCtx := httputil.UserFromRequest(r)

	var req UpdatePrivacyRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if err := h.validator.Validate(req); err != nil {
		httputil.ValidationError(w, err)
		return
	}

	user, err := h.service.UpdatePrivacy(
		r.Context(),
		userCtx.ID,
		req.UsernameSearchable,
		model.DMPolicy(req.DMPolicy),
	)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, UpdatePrivacyResponse{
		UsernameSearchable: user.UsernameSearchable,
		DMPolicy:           string(user.DMPolicy),
	})
}
