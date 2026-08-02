package user

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Service) SearchUsers(
	ctx context.Context,
	viewerID uuid.UUID,
	rawQuery string,
	limit int32,
) ([]model.UserSearchResult, error) {
	query := NormalizeUsername(rawQuery)
	if len(query) < searchMinQueryLength {
		return []model.UserSearchResult{}, nil
	}

	if limit <= 0 {
		limit = searchDefaultLimit
	}
	if limit > searchMaxLimit {
		limit = searchMaxLimit
	}

	return s.repo.SearchByUsernamePrefix(ctx, query, viewerID, limit)
}

func (s *Service) GetPublicProfile(
	ctx context.Context,
	viewerID uuid.UUID,
	rawUsername string,
) (model.PublicUserProfile, error) {
	username := NormalizeUsername(rawUsername)
	if err := ValidateUsername(username); err != nil {
		return model.PublicUserProfile{}, ErrUserNotDiscoverable
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.PublicUserProfile{}, ErrUserNotDiscoverable
		}
		return model.PublicUserProfile{}, err
	}

	if !user.OnboardingComplete() || !user.UsernameSearchable {
		return model.PublicUserProfile{}, ErrUserNotDiscoverable
	}

	canMessage, err := s.canViewerMessageUser(ctx, viewerID, user)
	if err != nil {
		return model.PublicUserProfile{}, err
	}

	return model.PublicUserProfile{
		ID:         user.ID,
		Username:   user.Username,
		AvatarURL:  user.AvatarURL,
		CanMessage: canMessage,
	}, nil
}

func (s *Service) ChangeUsername(ctx context.Context, userID uuid.UUID, rawUsername string) (model.User, error) {
	username := NormalizeUsername(rawUsername)
	if err := ValidateUsername(username); err != nil {
		return model.User{}, err
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	if user.Username == username {
		return user, nil
	}

	exists, err := s.repo.CheckUsernameExists(ctx, username)
	if err != nil {
		return model.User{}, err
	}
	if exists {
		return model.User{}, ErrUsernameTaken
	}

	if err := s.repo.UpdateUsername(ctx, userID, username); err != nil {
		if errors.Is(err, repository.ErrUniqueAlreadyExists) {
			return model.User{}, ErrUsernameTaken
		}
		return model.User{}, err
	}

	user.Username = username
	return user, nil
}

func (s *Service) UpdatePrivacy(
	ctx context.Context,
	userID uuid.UUID,
	searchable bool,
	policy model.DMPolicy,
) (model.User, error) {
	if policy != model.DMPolicyEveryone && policy != model.DMPolicyNobody {
		return model.User{}, errors.New("invalid dm policy")
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	if err := s.repo.UpdatePrivacy(ctx, userID, searchable, policy); err != nil {
		return model.User{}, err
	}

	user.UsernameSearchable = searchable
	user.DMPolicy = policy
	return user, nil
}

func (s *Service) ResolveUsername(ctx context.Context, rawUsername string) (model.User, error) {
	username := NormalizeUsername(rawUsername)
	if err := ValidateUsername(username); err != nil {
		return model.User{}, ErrUserNotDiscoverable
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.User{}, ErrUserNotDiscoverable
		}
		return model.User{}, err
	}

	if !user.OnboardingComplete() || !user.UsernameSearchable {
		return model.User{}, ErrUserNotDiscoverable
	}

	return user, nil
}

func (s *Service) CanViewerMessageUser(ctx context.Context, viewerID uuid.UUID, target model.User) (bool, error) {
	return s.canViewerMessageUser(ctx, viewerID, target)
}

func (s *Service) canViewerMessageUser(ctx context.Context, viewerID uuid.UUID, target model.User) (bool, error) {
	if viewerID == target.ID {
		return false, nil
	}

	if target.DMPolicy.AllowsNewMessages() {
		return true, nil
	}

	// Existing conversations stay reachable even when dm_policy is nobody.
	hasConversation, err := s.hasExistingDM(ctx, viewerID, target.ID)
	if err != nil {
		return false, err
	}

	return hasConversation, nil
}

func (s *Service) hasExistingDM(ctx context.Context, viewerID, targetID uuid.UUID) (bool, error) {
	if s.conversations == nil {
		return false, nil
	}

	low, high := model.DMPairKey(viewerID, targetID)
	_, err := s.conversations.FindDM(ctx, low, high)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, err
}
