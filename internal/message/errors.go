package message

import "errors"

var (
	ErrNotMessageSender   = errors.New("only the sender can delete a message for everyone")
	ErrEmptyMessage       = errors.New("message content or attachments required")
	ErrInvalidAttachments = errors.New("invalid attachments")
	ErrAttachmentTooLarge = errors.New("attachment too large")
	ErrInvalidMediaKind   = errors.New("invalid media kind")
	ErrEmptyMessageIDs    = errors.New("message ids required")
	ErrTooManyMessages    = errors.New("too many messages")
)
