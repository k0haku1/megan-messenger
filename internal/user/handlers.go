package user

import (
	"errors"
	"megan-messenger/internal/httputil"
	"net/http"
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

// CurrentUser returns the current user
//
//	@Summary	Get current user
//	@Tags		user
//	@Produce	json
//	@Security	BearerAuth
//	@Success	200	{object}	MeResponse
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/users/me [get]
func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	userCtx := httputil.UserFromRequest(r)

	user, err := h.service.GetUser(r.Context(), userCtx.ID)
	if err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.SuccessData(w, MeResponse{
		ID:                 user.ID.String(),
		Phone:              user.Phone,
		Username:           user.Username,
		AvatarURL:          user.AvatarURL,
		HasPassword:        user.HasPassword(),
		OnboardingComplete: user.OnboardingComplete(),
	})
}

// UploadAvatar uploads a new avatar for the user
//
//	@Summary	Upload user avatar
//	@Tags		user
//	@Accept		multipart/form-data
//	@Produce	json
//	@Security	BearerAuth
//	@Param		avatar	formData	file	true	"Avatar file"
//	@Success	204
//	@Failure	400	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/users/me/avatar [post]
func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		httputil.ValidationError(w, httputil.FieldErrors{"avatar": "file too big"})
		return
	}

	file, _, err := r.FormFile("avatar")
	if err != nil {
		httputil.ValidationError(w, httputil.FieldErrors{"avatar": "failed to read file"})
		return
	}
	defer file.Close()

	filename, err := h.service.UploadAvatar(file)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidImage), errors.Is(err, ErrUploadAvatar):
			httputil.ValidationError(w, httputil.FieldErrors{"avatar": err.Error()})
			return
		}

		httputil.InternalError(w, r, err)
		return
	}

	ctx := r.Context()
	if err := h.service.UpdateAvatar(ctx, httputil.UserFromRequest(r).ID, filename); err != nil {
		httputil.InternalError(w, r, err)
		return
	}

	httputil.Success(w)
}
