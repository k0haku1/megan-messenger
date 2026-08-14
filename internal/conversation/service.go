package conversation

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"megan-messenger/internal/storage"
	"megan-messenger/internal/user"
	"time"

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
	urls     *storage.URLResolver
	folders  repository.FolderRepository
}

func NewService(
	repo repository.ConversationRepository,
	users repository.UserRepository,
	messages repository.MessageRepository,
	urls *storage.URLResolver,
	folders repository.FolderRepository,
) *Service {
	return &Service{repo: repo, users: users, messages: messages, urls: urls, folders: folders}
}

func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Conversation, error) {
	conversations, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range conversations {
		s.resolvePeerAvatar(ctx, conversations[i].Peer)
		conversations[i].FolderIDs = []uuid.UUID{}
	}

	if s.folders != nil {
		folderMap, err := s.folders.GetFolderIDsByConversations(ctx, userID)
		if err == nil {
			for i := range conversations {
				if ids, ok := folderMap[conversations[i].ID]; ok {
					conversations[i].FolderIDs = ids
				}
			}
		}
	}

	return conversations, nil
}

func (s *Service) resolvePeerAvatar(ctx context.Context, peer *model.ConversationPeer) {
	if peer == nil || s.urls == nil {
		return
	}
	peer.AvatarURL = s.urls.Resolve(ctx, peer.AvatarURL)
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
		if _, err := s.users.GetByID(ctx, memberID); err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				return model.Conversation{}, ErrInvalidGroupMember
			}
			return model.Conversation{}, err
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
		conv.Peer = s.peerToConversationPeer(ctx, peer)
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

	created.Peer = s.peerToConversationPeer(ctx, peer)

	return created, nil
}

func (s *Service) peerToConversationPeer(ctx context.Context, peer model.User) *model.ConversationPeer {
	avatarURL := peer.AvatarURL
	if s.urls != nil {
		avatarURL = s.urls.Resolve(ctx, peer.AvatarURL)
	}
	return &model.ConversationPeer{
		ID:        peer.ID,
		Username:  peer.Username,
		AvatarURL: avatarURL,
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

	_, _, _ = s.repo.AdvanceLastReadAt(ctx, conv.ID, selfID, created.CreatedAt)

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

func (s *Service) LeaveGroup(ctx context.Context, userID, conversationID uuid.UUID) error {
	conv, err := s.repo.GetByID(ctx, conversationID)
	if err != nil {
		return err
	}

	if conv.Type != model.ConversationTypeGroup {
		return ErrNotGroupConversation
	}

	isMember, err := s.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return repository.ErrConversationNotFound
	}

	return s.repo.RemoveMember(ctx, conversationID, userID)
}

func (s *Service) MarkRead(ctx context.Context, userID, conversationID, messageID uuid.UUID) (time.Time, bool, error) {
	ok, err := s.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return time.Time{}, false, err
	}
	if !ok {
		return time.Time{}, false, repository.ErrConversationNotFound
	}

	message, err := s.messages.GetByID(ctx, messageID)
	if err != nil {
		return time.Time{}, false, err
	}
	if message.ConversationID != conversationID || message.DeletedAt != nil {
		return time.Time{}, false, repository.ErrMessageNotFound
	}

	lastReadAt, advanced, err := s.repo.AdvanceLastReadAt(ctx, conversationID, userID, message.CreatedAt)
	if err != nil {
		return time.Time{}, false, err
	}
	return lastReadAt, advanced, nil
}

func (s *Service) GetReadState(ctx context.Context, userID, conversationID uuid.UUID) (*time.Time, *time.Time, error) {
	ok, err := s.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, repository.ErrConversationNotFound
	}
	mine, err := s.repo.GetLastReadAt(ctx, conversationID, userID)
	if err != nil {
		return nil, nil, err
	}
	others, err := s.repo.GetOthersReadWatermark(ctx, conversationID, userID)
	if err != nil {
		return nil, nil, err
	}
	return mine, others, nil
}

func (s *Service) GetOthersReadWatermark(ctx context.Context, userID, conversationID uuid.UUID) (*time.Time, error) {
	ok, err := s.repo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, repository.ErrConversationNotFound
	}
	return s.repo.GetOthersReadWatermark(ctx, conversationID, userID)
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
