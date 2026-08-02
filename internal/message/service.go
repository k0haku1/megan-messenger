package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"
	"megan-messenger/internal/repository"
	"strconv"

	"github.com/google/uuid"
)

type Service struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	users            repository.UserRepository
}

func NewService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	users repository.UserRepository,
) *Service {
	return &Service{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		users:            users,
	}
}

func (s *Service) ListMessages(ctx context.Context, userID, conversationID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error) {
	ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, repository.ErrConversationNotFound
	}

	return s.messageRepo.ListMessages(ctx, conversationID, limit, cursor)
}

func (s *Service) SendMessage(
	ctx context.Context,
	userID, conversationID uuid.UUID,
	content string,
) (model.Message, error) {
	ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return model.Message{}, err
	}
	if !ok {
		return model.Message{}, repository.ErrConversationNotFound
	}

	sender, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return model.Message{}, err
	}

	msg, err := model.NewMessage(conversationID, content, sender)
	if err != nil {
		return model.Message{}, err
	}

	return s.messageRepo.CreateMessage(ctx, msg)
}

func (s *Service) GenerateCursor(message model.Message) string {
	c := pagination.Cursor{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
	}

	b, _ := json.Marshal(c)

	return base64.StdEncoding.EncodeToString(b)
}

func (s *Service) ParseCursor(cursorEncoded string) (pagination.Cursor, error) {
	var cursor pagination.Cursor

	decoded, err := base64.StdEncoding.DecodeString(cursorEncoded)
	if err != nil {
		return cursor, err
	}

	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return cursor, err
	}

	return cursor, nil
}

func normalizeLimit(limit string, max int, fallback int) int {
	if limit == "" {
		return fallback
	}

	result, err := strconv.Atoi(limit)
	if err != nil {
		return fallback
	}

	if result < 0 || result > max {
		return fallback
	}

	return result
}
