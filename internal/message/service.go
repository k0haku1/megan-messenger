package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"
	"megan-messenger/internal/repository"
	"strconv"
	"unicode/utf8"

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

	return s.messageRepo.ListMessages(ctx, conversationID, userID, limit, cursor)
}

func (s *Service) SendMessage(
	ctx context.Context,
	userID, conversationID uuid.UUID,
	content string,
	replyToID *uuid.UUID,
) (model.Message, error) {
	ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return model.Message{}, err
	}
	if !ok {
		return model.Message{}, repository.ErrConversationNotFound
	}
	if replyToID != nil {
		replied, err := s.messageRepo.GetByID(ctx, *replyToID)
		if err != nil {
			return model.Message{}, err
		}
		if replied.ConversationID != conversationID {
			return model.Message{}, repository.ErrMessageNotFound
		}
	}

	sender, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return model.Message{}, err
	}

	msg, err := model.NewMessage(conversationID, content, sender)
	if err != nil {
		return model.Message{}, err
	}
	msg.ReplyToID = replyToID

	return s.messageRepo.CreateMessage(ctx, msg)
}

func (s *Service) HideForUser(ctx context.Context, userID, messageID uuid.UUID) error {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	ok, err := s.conversationRepo.IsMember(ctx, message.ConversationID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return repository.ErrConversationNotFound
	}
	return s.messageRepo.HideForUser(ctx, messageID, userID)
}

func (s *Service) DeleteForEveryone(ctx context.Context, userID, messageID uuid.UUID) (model.Message, error) {
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return model.Message{}, err
	}
	ok, err := s.conversationRepo.IsMember(ctx, message.ConversationID, userID)
	if err != nil {
		return model.Message{}, err
	}
	if !ok {
		return model.Message{}, repository.ErrConversationNotFound
	}
	if message.Sender.ID != userID {
		return model.Message{}, ErrNotMessageSender
	}
	if _, err := s.messageRepo.DeleteForEveryone(ctx, messageID, userID); err != nil {
		return model.Message{}, err
	}
	updated, err := s.messageRepo.GetByIDForUser(ctx, messageID, userID)
	if err != nil {
		return model.Message{}, err
	}
	return updated, nil
}

func (s *Service) ForwardMessage(ctx context.Context, userID, messageID, targetConversationID uuid.UUID) (model.Message, error) {
	source, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return model.Message{}, err
	}
	if source.DeletedAt != nil {
		return model.Message{}, repository.ErrMessageNotFound
	}
	for _, conversationID := range []uuid.UUID{source.ConversationID, targetConversationID} {
		ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
		if err != nil {
			return model.Message{}, err
		}
		if !ok {
			return model.Message{}, repository.ErrConversationNotFound
		}
	}
	sender, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return model.Message{}, err
	}
	message, err := model.NewMessage(targetConversationID, source.Content, sender)
	if err != nil {
		return model.Message{}, err
	}
	message.ForwardedFromID = &source.ID
	return s.messageRepo.CreateMessage(ctx, message)
}

func (s *Service) SetReaction(ctx context.Context, userID, messageID uuid.UUID, emoji string, add bool) (model.Message, error) {
	if !utf8.ValidString(emoji) || utf8.RuneCountInString(emoji) > 8 {
		return model.Message{}, model.ErrInvalidMessageContent
	}
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return model.Message{}, err
	}
	if message.DeletedAt != nil {
		return model.Message{}, repository.ErrMessageNotFound
	}
	ok, err := s.conversationRepo.IsMember(ctx, message.ConversationID, userID)
	if err != nil {
		return model.Message{}, err
	}
	if !ok {
		return model.Message{}, repository.ErrConversationNotFound
	}
	if add {
		err = s.messageRepo.AddReaction(ctx, messageID, userID, emoji)
	} else {
		err = s.messageRepo.RemoveReaction(ctx, messageID, userID, emoji)
	}
	if err != nil {
		return model.Message{}, err
	}
	updated, err := s.messageRepo.GetByIDForUser(ctx, messageID, userID)
	if err != nil {
		return model.Message{}, err
	}
	updated.ReactionUpdatedBy = &userID
	return updated, nil
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
