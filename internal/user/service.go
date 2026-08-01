package user

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log/slog"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type Service struct {
	repo             repository.UserRepository
	avatarsUploadDir string
}

func NewService(repo repository.UserRepository, avatarsUploadDir string) *Service {
	return &Service{
		repo:             repo,
		avatarsUploadDir: avatarsUploadDir,
	}
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateAvatar(ctx context.Context, id uuid.UUID, url string) error {
	return s.repo.ChangeAvatar(ctx, id, url)
}

func (s *Service) UploadAvatar(file multipart.File) (string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	if format != "jpg" && format != "jpeg" && format != "png" && format != "webp" {
		return "", ErrInvalidImage
	}
	dstImg := imaging.Fill(img, 128, 128, imaging.Center, imaging.Lanczos)

	resultFilename := fmt.Sprintf("%s.%s", uuid.New().String(), format)
	filePath := filepath.Join(s.avatarsUploadDir, resultFilename)
	out, err := os.Create(filePath)
	if err != nil {
		slog.Error("file upload", "err", err, "dir", s.avatarsUploadDir)
		return "", ErrUploadAvatar
	}
	defer out.Close()

	if err := jpeg.Encode(out, dstImg, &jpeg.Options{Quality: 80}); err != nil {
		return "", err
	}

	return resultFilename, nil
}
