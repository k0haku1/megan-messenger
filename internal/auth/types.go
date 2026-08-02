package auth

type StartPhoneRequest struct {
	Phone string `json:"phone" validate:"required"`
}

type VerifyPhoneRequest struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

type VerifyPasswordChallengeRequest struct {
	ChallengeToken string `json:"challengeToken" validate:"required"`
	Password       string `json:"password" validate:"required,min=4,max=72"`
}

type SetUsernameRequest struct {
	Username string `json:"username" validate:"required,min=5,max=32"`
}

type SetPasswordRequest struct {
	Password        string `json:"password" validate:"required,min=4,max=72"`
	CurrentPassword string `json:"currentPassword,omitempty"`
}

type RemovePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=4,max=72"`
}

type AuthResult struct {
	AccessToken     string `json:"accessToken,omitempty"`
	RefreshToken    string `json:"refreshToken,omitempty"`
	NeedUsername    bool   `json:"needUsername"`
	NeedPassword    bool   `json:"needPassword"`
	ChallengeToken  string `json:"challengeToken,omitempty"`
}

type CompleteUsernameResponse struct {
	AuthResult
	User map[string]any `json:"user"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
