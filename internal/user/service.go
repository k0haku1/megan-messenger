package user

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"megan-messenger/internal/storage"
	"mime/multipart"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

type Service struct {
	repo          repository.UserRepository
	conversations repository.ConversationRepository
	store         storage.ObjectStore
	urls          *storage.URLResolver
}

func NewService(
	repo repository.UserRepository,
	conversations repository.ConversationRepository,
	store storage.ObjectStore,
	urls *storage.URLResolver,
) *Service {
	return &Service{
		repo:          repo,
		conversations: conversations,
		store:         store,
		urls:          urls,
	}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return model.User{}, err
	}
	s.resolveUserAvatar(ctx, &user)
	return user, nil
}

func (s *Service) resolveUserAvatar(ctx context.Context, user *model.User) {
	if user == nil || s.urls == nil {
		return
	}
	user.AvatarURL = s.urls.Resolve(ctx, user.AvatarURL)
}

func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File) (string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return "", ErrInvalidImage
	}
	if format != "jpg" && format != "jpeg" && format != "png" && format != "webp" && format != "gif" {
		return "", ErrInvalidImage
	}

	dstImg := imaging.Fill(img, 256, 256, imaging.Center, imaging.Lanczos)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dstImg, &jpeg.Options{Quality: 85}); err != nil {
		return "", err
	}

	current, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	key := storage.AvatarKey(userID, "jpg")
	if err := s.store.Put(ctx, key, bytes.NewReader(buf.Bytes()), int64(buf.Len()), "image/jpeg"); err != nil {
		return "", fmt.Errorf("%w: %v", ErrUploadAvatar, err)
	}

	if err := s.repo.ChangeAvatar(ctx, userID, key); err != nil {
		_ = s.store.Delete(ctx, key)
		return "", err
	}

	if current.AvatarURL != "" && containsSlash(current.AvatarURL) {
		_ = s.store.Delete(ctx, current.AvatarURL)
	}

	if s.urls == nil {
		return key, nil
	}
	return s.urls.Resolve(ctx, key), nil
}

func containsSlash(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			return true
		}
	}
	return false
}
