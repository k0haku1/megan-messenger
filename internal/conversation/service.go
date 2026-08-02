package conversation

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"megan-messenger/internal/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrCannotDMYourself     = errors.New("cannot dm yourself")
	ErrCannotMessageUser    = user.ErrCannotMessageUser
	ErrConversationNotFound = repository.ErrConversationNotFound
)

type Service struct {
	repo     repository.ConversationRepository
	users    repository.UserRepository
	messages repository.MessageRepository
}

func NewService(
	repo repository.ConversationRepository,
	users repository.UserRepository,
	messages repository.MessageRepository,
) *Service {
	return &Service{repo: repo, users: users, messages: messages}
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *Service) CreateGroup(ctx context.Context, creatorID uuid.UUID, title string, memberIDs []uuid.UUID) (model.Conversation, error) {
	conv, err := model.NewGroupConversation(title)
	if err != nil {
		return model.Conversation{}, err
	}

	created, err := s.repo.Create(ctx, conv)
	if err != nil {
		return model.Conversation{}, err
	}

	if err := s.repo.AddMember(ctx, created.ID, creatorID); err != nil {
		return model.Conversation{}, err
	}

	for _, memberID := range memberIDs {
		if memberID == creatorID {
			continue
		}
		if err := s.repo.AddMember(ctx, created.ID, memberID); err != nil {
			return model.Conversation{}, err
		}
	}

	return created, nil
}

func (s *Service) GetOrCreateDM(ctx context.Context, selfID, peerID uuid.UUID) (model.Conversation, error) {
	if selfID == peerID {
		return model.Conversation{}, ErrCannotDMYourself
	}

	peer, err := s.users.GetByID(ctx, peerID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.Conversation{}, user.ErrUserNotDiscoverable
		}
		return model.Conversation{}, err
	}

	if !peer.OnboardingComplete() {
		return model.Conversation{}, user.ErrUserNotDiscoverable
	}

	low, high := model.DMPairKey(selfID, peerID)

	convID, err := s.repo.FindDM(ctx, low, high)
	if err == nil {
		conv, err := s.repo.GetByID(ctx, convID)
		if err != nil {
			return model.Conversation{}, err
		}
		conv.Peer = peerToConversationPeer(peer)
		return conv, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.Conversation{}, err
	}

	canMessage, err := s.usersCanMessage(ctx, selfID, peer)
	if err != nil {
		return model.Conversation{}, err
	}
	if !canMessage {
		return model.Conversation{}, ErrCannotMessageUser
	}

	conv := model.NewDMConversation()
	created, err := s.repo.Create(ctx, conv)
	if err != nil {
		return model.Conversation{}, err
	}

	if err := s.repo.AddMember(ctx, created.ID, selfID); err != nil {
		return model.Conversation{}, err
	}
	if err := s.repo.AddMember(ctx, created.ID, peerID); err != nil {
		return model.Conversation{}, err
	}
	if err := s.repo.CreateDMPair(ctx, low, high, created.ID); err != nil {
		return model.Conversation{}, err
	}

	created.Peer = peerToConversationPeer(peer)

	return created, nil
}

func peerToConversationPeer(peer model.User) *model.ConversationPeer {
	return &model.ConversationPeer{
		ID:        peer.ID,
		Username:  peer.Username,
		AvatarURL: peer.AvatarURL,
	}
}

func (s *Service) GetOrCreateDMByUsername(ctx context.Context, selfID uuid.UUID, rawUsername string) (model.Conversation, error) {
	peer, err := s.resolveDiscoverableUser(ctx, rawUsername)
	if err != nil {
		return model.Conversation{}, err
	}

	return s.GetOrCreateDM(ctx, selfID, peer.ID)
}

func (s *Service) SendDMMessage(
	ctx context.Context,
	selfID uuid.UUID,
	userID *uuid.UUID,
	username *string,
	content string,
) (model.Conversation, model.Message, error) {
	var conv model.Conversation
	var err error

	if userID != nil {
		conv, err = s.GetOrCreateDM(ctx, selfID, *userID)
	} else {
		conv, err = s.GetOrCreateDMByUsername(ctx, selfID, *username)
	}
	if err != nil {
		return model.Conversation{}, model.Message{}, err
	}

	sender, err := s.users.GetByID(ctx, selfID)
	if err != nil {
		return model.Conversation{}, model.Message{}, err
	}

	msg, err := model.NewMessage(conv.ID, content, sender)
	if err != nil {
		return model.Conversation{}, model.Message{}, err
	}

	created, err := s.messages.CreateMessage(ctx, msg)
	if err != nil {
		return model.Conversation{}, model.Message{}, err
	}

	return conv, created, nil
}

func (s *Service) resolveDiscoverableUser(ctx context.Context, rawUsername string) (model.User, error) {
	username := user.NormalizeUsername(rawUsername)
	if err := user.ValidateUsername(username); err != nil {
		return model.User{}, user.ErrUserNotDiscoverable
	}

	target, err := s.users.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return model.User{}, user.ErrUserNotDiscoverable
		}
		return model.User{}, err
	}

	if !target.OnboardingComplete() || !target.UsernameSearchable {
		return model.User{}, user.ErrUserNotDiscoverable
	}

	return target, nil
}

func (s *Service) usersCanMessage(ctx context.Context, viewerID uuid.UUID, target model.User) (bool, error) {
	if target.DMPolicy.AllowsNewMessages() {
		return true, nil
	}

	low, high := model.DMPairKey(viewerID, target.ID)
	_, err := s.repo.FindDM(ctx, low, high)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func (s *Service) JoinBySlug(ctx context.Context, userID uuid.UUID, slug string) (model.Conversation, error) {
	conv, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return model.Conversation{}, err
	}

	if conv.Type != model.ConversationTypeGroup {
		return model.Conversation{}, repository.ErrConversationNotFound
	}

	if err := s.repo.AddMember(ctx, conv.ID, userID); err != nil {
		return model.Conversation{}, err
	}

	return conv, nil
}

func (s *Service) EnsureMember(ctx context.Context, userID, conversationID uuid.UUID) (model.Conversation, error) {
	ok, err := s.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return model.Conversation{}, err
	}
	if !ok {
		return model.Conversation{}, repository.ErrConversationNotFound
	}
	return s.repo.GetByID(ctx, conversationID)
}
