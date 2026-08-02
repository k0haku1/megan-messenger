package auth

import (
	"errors"
	"megan-messenger/internal/httputil"
	userpkg "megan-messenger/internal/user"
	"megan-messenger/internal/model"
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

// StartPhoneAuth sends an OTP to the phone number
//
//	@Summary	Start phone auth (send OTP)
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		input	body	StartPhoneRequest	true	"Phone"
//	@Success	204
//	@Failure	400	{object}	httputil.ErrorResponse
//	@Failure	429	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/auth/phone/start [post]
func (h *Handler) StartPhoneAuth(w http.ResponseWriter, r *http.Request) {
	var req StartPhoneRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fieldErrs := h.validator.Validate(&req); fieldErrs != nil {
		httputil.ValidationError(w, fieldErrs)
		return
	}

	if err := h.service.StartPhoneAuth(r.Context(), req.Phone); err != nil {
		switch {
		case errors.Is(err, ErrInvalidPhone):
			httputil.ValidationError(w, map[string]string{"phone": err.Error()})
		case errors.Is(err, ErrOTPResendCooldown):
			httputil.SimpleError(w, http.StatusTooManyRequests, "please wait before requesting another code")
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.NoContent(w)
}

// VerifyPhone verifies OTP and returns tokens or password challenge
//
//	@Summary	Verify phone OTP
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		input	body		VerifyPhoneRequest	true	"Phone and code"
//	@Success	200		{object}	AuthResult
//	@Failure	400		{object}	httputil.ErrorResponse
//	@Failure	429		{object}	httputil.ErrorResponse
//	@Failure	500		{object}	httputil.ErrorResponse
//	@Router		/auth/phone/verify [post]
func (h *Handler) VerifyPhone(w http.ResponseWriter, r *http.Request) {
	var req VerifyPhoneRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fieldErrs := h.validator.Validate(&req); fieldErrs != nil {
		httputil.ValidationError(w, fieldErrs)
		return
	}

	result, err := h.service.VerifyPhone(r.Context(), req.Phone, req.Code)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPhone):
			httputil.ValidationError(w, map[string]string{"phone": err.Error()})
		case errors.Is(err, ErrInvalidCode), errors.Is(err, ErrCodeExpired):
			httputil.ValidationError(w, map[string]string{"code": err.Error()})
		case errors.Is(err, ErrTooManyAttempts):
			httputil.SimpleError(w, http.StatusTooManyRequests, "too many attempts, request a new code")
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	if !result.NeedPassword {
		h.setRefreshTokenCookie(w, result.RefreshToken)
	}
	httputil.SuccessData(w, result)
}

// VerifyPasswordChallenge completes login when cloud password is set
//
//	@Summary	Verify cloud password challenge
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		input	body		VerifyPasswordChallengeRequest	true	"Challenge and password"
//	@Success	200		{object}	AuthResult
//	@Failure	401		{object}	httputil.ErrorResponse
//	@Failure	500		{object}	httputil.ErrorResponse
//	@Router		/auth/password/verify [post]
func (h *Handler) VerifyPasswordChallenge(w http.ResponseWriter, r *http.Request) {
	var req VerifyPasswordChallengeRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fieldErrs := h.validator.Validate(&req); fieldErrs != nil {
		httputil.ValidationError(w, fieldErrs)
		return
	}

	result, err := h.service.VerifyPasswordChallenge(r.Context(), req.ChallengeToken, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidChallenge), errors.Is(err, ErrInvalidCredentials):
			httputil.Unauthorized(w, ErrInvalidCredentials.Error())
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	h.setRefreshTokenCookie(w, result.RefreshToken)
	httputil.SuccessData(w, result)
}

// CompleteUsername sets username during onboarding
//
//	@Summary	Complete onboarding username
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		input	body		SetUsernameRequest	true	"Username"
//	@Success	200		{object}	CompleteUsernameResponse
//	@Failure	400		{object}	httputil.ErrorResponse
//	@Failure	401		{object}	httputil.ErrorResponse
//	@Failure	500		{object}	httputil.ErrorResponse
//	@Router		/auth/onboarding/username [post]
func (h *Handler) CompleteUsername(w http.ResponseWriter, r *http.Request) {
	var req SetUsernameRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fieldErrs := h.validator.Validate(&req); fieldErrs != nil {
		httputil.ValidationError(w, fieldErrs)
		return
	}

	userCtx := httputil.UserFromRequest(r)
	result, user, err := h.service.CompleteUsername(r.Context(), userCtx.ID, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, ErrUsernameExists):
			httputil.ValidationError(w, map[string]string{"username": err.Error()})
		case errors.Is(err, ErrUsernameAlreadySet):
			httputil.Conflict(w, err.Error())
		case errors.Is(err, userpkg.ErrInvalidUsername):
			httputil.ValidationError(w, map[string]string{"username": err.Error()})
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	h.setRefreshTokenCookie(w, result.RefreshToken)
	httputil.SuccessData(w, CompleteUsernameResponse{
		AuthResult: result,
		User:       toUserResponse(user),
	})
}

// SetPassword sets or changes cloud password
//
//	@Summary	Set or change cloud password
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		input	body	SetPasswordRequest	true	"Password"
//	@Success	204
//	@Failure	400	{object}	httputil.ErrorResponse
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/auth/password [post]
func (h *Handler) SetPassword(w http.ResponseWriter, r *http.Request) {
	var req SetPasswordRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fieldErrs := h.validator.Validate(&req); fieldErrs != nil {
		httputil.ValidationError(w, fieldErrs)
		return
	}

	userCtx := httputil.UserFromRequest(r)
	if err := h.service.SetPassword(r.Context(), userCtx.ID, req.Password, req.CurrentPassword); err != nil {
		switch {
		case errors.Is(err, ErrPasswordRequired), errors.Is(err, ErrInvalidPassword):
			httputil.ValidationError(w, map[string]string{"currentPassword": err.Error()})
		default:
			httputil.InternalError(w, r, err)
		}
		return
	}

	httputil.NoContent(w)
}

// RemovePassword removes cloud password
//
//	@Summary	Remove cloud password
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		input	body	RemovePasswordRequest	true	"Current password"
//	@Success	204
//	@Failure	400	{object}	httputil.ErrorResponse
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/auth/password [delete]
func (h *Handler) RemovePassword(w http.ResponseWriter, r *http.Request) {
	var req RemovePasswordRequest
	if err := httputil.Read(r, &req); err != nil {
		httputil.InvalidRequestBody(w)
		return
	}
	if fieldErrs := h.validator.Validate(&req); fieldErrs != nil {
		httputil.ValidationError(w, fieldErrs)
		return
	}

	userCtx := httputil.UserFromRequest(r)
	if err := h.service.RemovePassword(r.Context(), userCtx.ID, req.CurrentPassword); err != nil {
		if errors.Is(err, ErrInvalidPassword) {
			httputil.ValidationError(w, map[string]string{"currentPassword": err.Error()})
			return
		}
		httputil.InternalError(w, r, err)
		return
	}

	httputil.NoContent(w)
}

// Refresh refreshes the access token
//
//	@Summary	Refresh access token
//	@Tags		auth
//	@Produce	json
//	@Success	200	{object}	Tokens
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Router		/auth/refresh [post]
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httputil.Unauthorized(w, "no refresh token")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httputil.Unauthorized(w, "invalid refresh token")
		return
	}

	h.setRefreshTokenCookie(w, tokens.RefreshToken)
	httputil.SuccessData(w, tokens)
}

// Logout logs out a user
//
//	@Summary	Logout a user
//	@Tags		auth
//	@Produce	json
//	@Success	200
//	@Failure	401	{object}	httputil.ErrorResponse
//	@Failure	500	{object}	httputil.ErrorResponse
//	@Router		/auth/logout [post]
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httputil.Unauthorized(w, "no refresh token")
		return
	}

	if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
		httputil.InternalError(w, r, err)
		return
	}
	h.setRefreshTokenCookie(w, "")

	httputil.Success(w)
}

func (h *Handler) setRefreshTokenCookie(w http.ResponseWriter, refreshToken string) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/auth",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24 * 30,
	}
	if refreshToken == "" {
		cookie.MaxAge = -1
	}
	http.SetCookie(w, cookie)
}

func toUserResponse(user model.User) map[string]any {
	return map[string]any{
		"id":                 user.ID,
		"phone":              user.Phone,
		"username":           user.Username,
		"avatarUrl":          user.AvatarURL,
		"hasPassword":        user.HasPassword(),
		"onboardingComplete": user.OnboardingComplete(),
	}
}
