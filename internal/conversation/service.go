package conversation

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrCannotDMYourself     = errors.New("cannot dm yourself")
	ErrConversationNotFound = repository.ErrConversationNotFound
)

type Service struct {
	repo repository.ConversationRepository
}

func NewService(repo repository.ConversationRepository) *Service {
	return &Service{repo: repo}
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

	low, high := model.DMPairKey(selfID, peerID)

	convID, err := s.repo.FindDM(ctx, low, high)
	if err == nil {
		return s.repo.GetByID(ctx, convID)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return model.Conversation{}, err
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

	return created, nil
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
