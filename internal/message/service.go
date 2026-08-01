package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"
	"megan-messenger/internal/repository"
	"strconv"

	"github.com/google/uuid"
)

type Service struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
}

var (
	ErrConversationNotFound = errors.New("conversation not found")
)

func NewService(conversationRepo repository.ConversationRepository, messageRepo repository.MessageRepository) *Service {
	return &Service{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	}
}

func (s *Service) ListMessages(ctx context.Context, userID, conversationID uuid.UUID, limit int, cursor *pagination.Cursor) ([]model.Message, error) {
	ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrConversationNotFound
	}

	return s.messageRepo.ListMessages(ctx, conversationID, limit, cursor)
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
