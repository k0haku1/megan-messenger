package model

import (
	"time"

	"github.com/google/uuid"
)

type AttachmentKind string

const (
	AttachmentKindImage AttachmentKind = "image"
	AttachmentKindVideo AttachmentKind = "video"
	AttachmentKindFile  AttachmentKind = "file"
)

type MessageAttachment struct {
	ID             uuid.UUID      `json:"id"`
	MessageID      *uuid.UUID     `json:"messageId,omitempty"`
	ConversationID uuid.UUID      `json:"conversationId"`
	UploaderID     uuid.UUID      `json:"uploaderId"`
	ObjectKey      string         `json:"-"`
	ThumbKey       string         `json:"-"`
	Mime           string         `json:"mime"`
	Kind           AttachmentKind `json:"kind"`
	SizeBytes      int64          `json:"sizeBytes"`
	Width          *int           `json:"width,omitempty"`
	Height         *int           `json:"height,omitempty"`
	DurationMs     *int           `json:"durationMs,omitempty"`
	OriginalName   string         `json:"originalName"`
	URL            string         `json:"url,omitempty"`
	ThumbURL       string         `json:"thumbUrl,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
}
