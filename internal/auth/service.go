package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"megan-messenger/internal/config"
	"megan-messenger/internal/model"
	"megan-messenger/internal/notification"
	"megan-messenger/internal/repository"
	userpkg "megan-messenger/internal/user"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameExists     = errors.New("username already taken")
	ErrUsernameAlreadySet = errors.New("username already set")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidPhone       = errors.New("invalid phone number")
	ErrInvalidCode        = errors.New("invalid code")
	ErrCodeExpired        = errors.New("code expired")
	ErrTooManyAttempts    = errors.New("too many attempts")
	ErrOTPResendCooldown  = errors.New("otp resend cooldown")
	ErrInvalidChallenge   = errors.New("invalid challenge token")
	ErrPasswordRequired   = errors.New("current password required")
	ErrInvalidPassword    = errors.New("invalid password")
)

type Service struct {
	authenticator  *Authenticator
	userRepo       repository.UserRepository
	refreshRepo    repository.RefreshTokenRepository
	challengeRepo  repository.PasswordChallengeRepository
	smsSender      notification.SmsSender
	otpTTL         time.Duration
	otpResendAfter time.Duration
	otpMaxAttempts int
}

func NewService(
	authenticator *Authenticator,
	refreshRepo repository.RefreshTokenRepository,
	challengeRepo repository.PasswordChallengeRepository,
	userRepo repository.UserRepository,
	smsSender notification.SmsSender,
	authCfg config.AuthConfig,
) *Service {
	return &Service{
		authenticator:  authenticator,
		refreshRepo:    refreshRepo,
		challengeRepo:  challengeRepo,
		userRepo:       userRepo,
		smsSender:      smsSender,
		otpTTL:         authCfg.OTP.TTL,
		otpResendAfter: authCfg.OTP.ResendCooldown,
		otpMaxAttempts: authCfg.OTP.MaxAttempts,
	}
}

func (s *Service) StartPhoneAuth(ctx context.Context, rawPhone string) error {
	phone, err := NormalizePhone(rawPhone)
	if err != nil {
		return err
	}

	if existing, err := s.userRepo.GetPhoneOTP(ctx, phone); err == nil {
		if time.Since(existing.SentAt) < s.otpResendAfter {
			return ErrOTPResendCooldown
		}
	} else if !errors.Is(err, repository.ErrOTPNotFound) {
		return err
	}

	code := s.generateVerificationCode()
	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()
	if err := s.userRepo.UpsertPhoneOTP(ctx, phone, string(hashedCode), now.Add(s.otpTTL), now); err != nil {
		return err
	}

	return s.smsSender.SendCode(ctx, phone, code)
}

func (s *Service) VerifyPhone(ctx context.Context, rawPhone, code string) (AuthResult, error) {
	phone, err := NormalizePhone(rawPhone)
	if err != nil {
		return AuthResult{}, err
	}

	otp, err := s.userRepo.GetPhoneOTP(ctx, phone)
	if err != nil {
		if errors.Is(err, repository.ErrOTPNotFound) {
			return AuthResult{}, ErrInvalidCode
		}
		return AuthResult{}, err
	}

	if otp.Attempts >= s.otpMaxAttempts {
		return AuthResult{}, ErrTooManyAttempts
	}

	if time.Now().After(otp.ExpiresAt) {
		return AuthResult{}, ErrCodeExpired
	}

	if err := bcrypt.CompareHashAndPassword([]byte(otp.CodeHash), []byte(code)); err != nil {
		_ = s.userRepo.IncrementPhoneOTPAttempts(ctx, phone)
		return AuthResult{}, ErrInvalidCode
	}

	_ = s.userRepo.DeletePhoneOTP(ctx, phone)

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return AuthResult{}, err
		}
		created, createErr := s.userRepo.Create(ctx, model.NewUser(phone))
		if createErr != nil {
			return AuthResult{}, createErr
		}
		user = created
	}

	if user.HasPassword() {
		challengeToken, challengeErr := s.challengeRepo.Issue(ctx, user.ID)
		if challengeErr != nil {
			return AuthResult{}, challengeErr
		}
		return AuthResult{
			NeedPassword:   true,
			ChallengeToken: challengeToken,
		}, nil
	}

	return s.issueAuthResult(ctx, user)
}

func (s *Service) VerifyPasswordChallenge(ctx context.Context, challengeToken, password string) (AuthResult, error) {
	userID, err := s.challengeRepo.Consume(ctx, challengeToken)
	if err != nil {
		return AuthResult{}, ErrInvalidChallenge
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return AuthResult{}, err
	}

	if err := user.ComparePasswords(password); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	return s.issueAuthResult(ctx, user)
}

func (s *Service) CompleteUsername(ctx context.Context, userID uuid.UUID, rawUsername string) (AuthResult, model.User, error) {
	currentUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return AuthResult{}, model.User{}, err
	}
	if currentUser.Username != "" {
		return AuthResult{}, model.User{}, ErrUsernameAlreadySet
	}

	username := userpkg.NormalizeUsername(rawUsername)
	if err := userpkg.ValidateUsername(username); err != nil {
		return AuthResult{}, model.User{}, err
	}

	exists, err := s.userRepo.CheckUsernameExists(ctx, username)
	if err != nil {
		return AuthResult{}, model.User{}, err
	}
	if exists {
		return AuthResult{}, model.User{}, ErrUsernameExists
	}

	if err := s.userRepo.SetUsername(ctx, userID, username); err != nil {
		if errors.Is(err, repository.ErrUniqueAlreadyExists) {
			return AuthResult{}, model.User{}, ErrUsernameExists
		}
		return AuthResult{}, model.User{}, err
	}

	currentUser.Username = username
	result, err := s.issueAuthResult(ctx, currentUser)
	if err != nil {
		return AuthResult{}, model.User{}, err
	}
	return result, currentUser, nil
}

func (s *Service) SetPassword(ctx context.Context, userID uuid.UUID, password, currentPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.HasPassword() {
		if currentPassword == "" {
			return ErrPasswordRequired
		}
		if err := user.ComparePasswords(currentPassword); err != nil {
			return ErrInvalidPassword
		}
	}

	if err := user.ChangePassword(password); err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, user.ID, user.PasswordHash)
}

func (s *Service) RemovePassword(ctx context.Context, userID uuid.UUID, currentPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if !user.HasPassword() {
		return nil
	}
	if err := user.ComparePasswords(currentPassword); err != nil {
		return ErrInvalidPassword
	}
	return s.userRepo.ClearPassword(ctx, userID)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	return s.refreshRepo.Revoke(ctx, refreshToken)
}

func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.refreshRepo.RevokeAll(ctx, userID)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	userID, err := s.refreshRepo.Consume(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return Tokens{}, err
	}

	return s.issueTokens(ctx, user)
}

func (s *Service) issueAuthResult(ctx context.Context, user model.User) (AuthResult, error) {
	tokens, err := s.issueTokens(ctx, user)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		NeedUsername: !user.OnboardingComplete(),
		NeedPassword: false,
	}, nil
}

func (s *Service) issueTokens(ctx context.Context, user model.User) (Tokens, error) {
	claims := s.authenticator.GenerateClaims(user)

	accessToken, err := s.authenticator.GenerateToken(claims)
	if err != nil {
		return Tokens{}, err
	}
	refreshToken, err := s.refreshRepo.Issue(ctx, user.ID)
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) generateVerificationCode() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "000000"
	}
	return fmt.Sprintf("%06d", n)
}
