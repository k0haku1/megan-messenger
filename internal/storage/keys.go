package storage

import (
	"fmt"

	"github.com/google/uuid"
)

func AvatarKey(userID uuid.UUID, ext string) string {
	if ext == "" {
		ext = "jpg"
	}
	return fmt.Sprintf("avatars/%s/%s.%s", userID.String(), uuid.New().String(), ext)
}

func AttachmentKey(conversationID uuid.UUID, ext string) string {
	if ext == "" {
		ext = "bin"
	}
	return fmt.Sprintf("attachments/%s/%s.%s", conversationID.String(), uuid.New().String(), ext)
}

func ThumbKey(attachmentID uuid.UUID) string {
	return fmt.Sprintf("thumbs/%s.jpg", attachmentID.String())
}
