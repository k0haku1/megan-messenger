package folder

import (
	"context"
	"errors"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"

	"github.com/google/uuid"
)

const (
	maxFoldersPerUser = 20
	maxItemsPerFolder = 200
	minNameLen        = 1
	maxNameLen        = 64
)

var (
	ErrFolderNotFound    = errors.New("folder not found")
	ErrFolderLimitReached = errors.New("folder limit reached")
	ErrItemLimitReached  = errors.New("item limit reached")
	ErrInvalidName       = errors.New("invalid folder name")
	ErrNotMember         = errors.New("not a conversation member")
)

type Service struct {
	folders       repository.FolderRepository
	conversations repository.ConversationRepository
}

func NewService(folders repository.FolderRepository, conversations repository.ConversationRepository) *Service {
	return &Service{folders: folders, conversations: conversations}
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]model.ChatFolder, error) {
	return s.folders.ListByUser(ctx, userID)
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, name, icon string) (model.ChatFolder, error) {
	if len(name) < minNameLen || len(name) > maxNameLen {
		return model.ChatFolder{}, ErrInvalidName
	}

	count, err := s.folders.CountByUser(ctx, userID)
	if err != nil {
		return model.ChatFolder{}, err
	}
	if count >= maxFoldersPerUser {
		return model.ChatFolder{}, ErrFolderLimitReached
	}

	maxPos, err := s.folders.MaxPosition(ctx, userID)
	if err != nil {
		return model.ChatFolder{}, err
	}

	f := model.NewChatFolder(userID, name, icon, maxPos+1)
	return s.folders.Create(ctx, f)
}

func (s *Service) Update(ctx context.Context, userID, folderID uuid.UUID, name, icon string) (model.ChatFolder, error) {
	if len(name) < minNameLen || len(name) > maxNameLen {
		return model.ChatFolder{}, ErrInvalidName
	}

	f, err := s.folders.GetByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, repository.ErrFolderNotFound) {
			return model.ChatFolder{}, ErrFolderNotFound
		}
		return model.ChatFolder{}, err
	}
	if f.UserID != userID {
		return model.ChatFolder{}, ErrFolderNotFound
	}

	f.Name = name
	f.Icon = icon
	return s.folders.Update(ctx, f)
}

func (s *Service) Delete(ctx context.Context, userID, folderID uuid.UUID) error {
	f, err := s.folders.GetByID(ctx, folderID)
	if err != nil {
		if errors.Is(err, repository.ErrFolderNotFound) {
			return ErrFolderNotFound
		}
		return err
	}
	if f.UserID != userID {
		return ErrFolderNotFound
	}
	return s.folders.Delete(ctx, folderID, userID)
}

func (s *Service) Reorder(ctx context.Context, userID uuid.UUID, folderIDs []uuid.UUID) error {
	for i, id := range folderIDs {
		f, err := s.folders.GetByID(ctx, id)
		if err != nil || f.UserID != userID {
			continue
		}
		_ = s.folders.UpdatePosition(ctx, id, i, userID)
	}
	return nil
}

func (s *Service) AddItem(ctx context.Context, userID, folderID, conversationID uuid.UUID) error {
	f, err := s.folders.GetByID(ctx, folderID)
	if err != nil || f.UserID != userID {
		return ErrFolderNotFound
	}

	isMember, err := s.conversations.IsMember(ctx, conversationID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrNotMember
	}

	exists, err := s.folders.ItemExists(ctx, folderID, conversationID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	count, err := s.folders.CountItems(ctx, folderID)
	if err != nil {
		return err
	}
	if count >= maxItemsPerFolder {
		return ErrItemLimitReached
	}

	return s.folders.AddItem(ctx, folderID, conversationID)
}

func (s *Service) RemoveItem(ctx context.Context, userID, folderID, conversationID uuid.UUID) error {
	f, err := s.folders.GetByID(ctx, folderID)
	if err != nil || f.UserID != userID {
		return ErrFolderNotFound
	}
	return s.folders.RemoveItem(ctx, folderID, conversationID)
}
