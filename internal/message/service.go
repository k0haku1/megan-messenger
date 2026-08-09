package message

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"megan-messenger/internal/model"
	"megan-messenger/internal/pagination"
	"megan-messenger/internal/repository"
	"megan-messenger/internal/storage"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp"
)

const (
	maxImageBytes = 15 << 20
	maxVideoBytes = 50 << 20
	maxFileBytes  = 25 << 20
)

type Service struct {
	conversationRepo repository.ConversationRepository
	messageRepo      repository.MessageRepository
	attachmentRepo   repository.AttachmentRepository
	users            repository.UserRepository
	store            storage.ObjectStore
	urls             *storage.URLResolver
}

func NewService(
	conversationRepo repository.ConversationRepository,
	messageRepo repository.MessageRepository,
	attachmentRepo repository.AttachmentRepository,
	users repository.UserRepository,
	store storage.ObjectStore,
	urls *storage.URLResolver,
) *Service {
	return &Service{
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
		attachmentRepo:   attachmentRepo,
		users:            users,
		store:            store,
		urls:             urls,
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

	messages, err := s.messageRepo.ListMessages(ctx, conversationID, userID, limit, cursor)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateMessages(ctx, messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *Service) SendMessage(
	ctx context.Context,
	userID, conversationID uuid.UUID,
	content string,
	replyToID *uuid.UUID,
	attachmentIDs []uuid.UUID,
) (model.Message, error) {
	content = strings.TrimSpace(content)
	if content == "" && len(attachmentIDs) == 0 {
		return model.Message{}, ErrEmptyMessage
	}
	if utf8.RuneCountInString(content) > 5000 {
		return model.Message{}, model.ErrInvalidMessageContent
	}

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

	if len(attachmentIDs) > 0 {
		count, err := s.attachmentRepo.CountPending(ctx, userID, conversationID, attachmentIDs)
		if err != nil {
			return model.Message{}, err
		}
		if count != int64(len(attachmentIDs)) {
			return model.Message{}, ErrInvalidAttachments
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

	created, err := s.messageRepo.CreateMessage(ctx, msg)
	if err != nil {
		return model.Message{}, err
	}

	if len(attachmentIDs) > 0 {
		if err := s.attachmentRepo.AttachToMessage(ctx, created.ID, userID, conversationID, attachmentIDs); err != nil {
			return model.Message{}, err
		}
	}

	hydrated := []model.Message{created}
	if err := s.hydrateMessages(ctx, hydrated); err != nil {
		return model.Message{}, err
	}
	return hydrated[0], nil
}

func (s *Service) UploadAttachment(
	ctx context.Context,
	userID, conversationID uuid.UUID,
	filename string,
	reader io.Reader,
	size int64,
) (model.MessageAttachment, error) {
	ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return model.MessageAttachment{}, err
	}
	if !ok {
		return model.MessageAttachment{}, repository.ErrConversationNotFound
	}

	limited := io.LimitReader(reader, maxVideoBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return model.MessageAttachment{}, err
	}
	if int64(len(data)) != size && size > 0 && int64(len(data)) < size {
		// size hint may be unknown; use actual
	}
	actualSize := int64(len(data))

	detected := mimetype.Detect(data)
	mime := detected.String()
	kind, ext, maxSize := classifyMime(mime, filename)
	if actualSize > maxSize {
		return model.MessageAttachment{}, ErrAttachmentTooLarge
	}

	id := uuid.Must(uuid.NewV7())
	objectKey := storage.AttachmentKey(conversationID, ext)
	if err := s.store.Put(ctx, objectKey, bytes.NewReader(data), actualSize, mime); err != nil {
		return model.MessageAttachment{}, fmt.Errorf("store attachment: %w", err)
	}

	var width, height *int
	var thumbKey string
	if kind == model.AttachmentKindImage {
		w, h, thumb, err := makeImageThumb(data)
		if err == nil {
			width, height = &w, &h
			thumbKey = storage.ThumbKey(id)
			if err := s.store.Put(ctx, thumbKey, bytes.NewReader(thumb), int64(len(thumb)), "image/jpeg"); err != nil {
				thumbKey = ""
			}
		}
	}

	att := model.MessageAttachment{
		ID:             id,
		ConversationID: conversationID,
		UploaderID:     userID,
		ObjectKey:      objectKey,
		ThumbKey:       thumbKey,
		Mime:           mime,
		Kind:           kind,
		SizeBytes:      actualSize,
		Width:          width,
		Height:         height,
		OriginalName:   filepath.Base(filename),
		CreatedAt:      time.Now().UTC(),
	}

	created, err := s.attachmentRepo.Create(ctx, att)
	if err != nil {
		_ = s.store.Delete(ctx, objectKey)
		if thumbKey != "" {
			_ = s.store.Delete(ctx, thumbKey)
		}
		return model.MessageAttachment{}, err
	}
	s.resolveAttachment(ctx, &created)
	return created, nil
}

func (s *Service) ListMedia(
	ctx context.Context,
	userID, conversationID uuid.UUID,
	kind string,
	limit int,
	cursor *pagination.Cursor,
) ([]model.MessageAttachment, error) {
	ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, repository.ErrConversationNotFound
	}
	if kind != "" && kind != "image" && kind != "video" && kind != "file" {
		return nil, ErrInvalidMediaKind
	}
	items, err := s.attachmentRepo.ListConversationMedia(ctx, conversationID, kind, limit, cursor)
	if err != nil {
		return nil, err
	}
	for i := range items {
		s.resolveAttachment(ctx, &items[i])
		// Gallery prefers thumbs; keep original URL for open-on-demand.
	}
	return items, nil
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
	hydrated := []model.Message{updated}
	_ = s.hydrateMessages(ctx, hydrated)
	return hydrated[0], nil
}

const maxForwardMessages = 50

func (s *Service) ForwardMessage(ctx context.Context, userID, messageID, targetConversationID uuid.UUID) (model.Message, error) {
	source, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return model.Message{}, err
	}
	messages, err := s.ForwardMessages(ctx, userID, source.ConversationID, targetConversationID, []uuid.UUID{messageID})
	if err != nil {
		return model.Message{}, err
	}
	if len(messages) == 0 {
		return model.Message{}, repository.ErrMessageNotFound
	}
	return messages[0], nil
}

func (s *Service) ForwardMessages(
	ctx context.Context,
	userID, sourceConversationID, targetConversationID uuid.UUID,
	messageIDs []uuid.UUID,
) ([]model.Message, error) {
	if len(messageIDs) == 0 {
		return nil, ErrEmptyMessageIDs
	}
	if len(messageIDs) > maxForwardMessages {
		return nil, ErrTooManyMessages
	}

	for _, conversationID := range []uuid.UUID{sourceConversationID, targetConversationID} {
		ok, err := s.conversationRepo.IsMember(ctx, conversationID, userID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, repository.ErrConversationNotFound
		}
	}

	seen := make(map[uuid.UUID]struct{}, len(messageIDs))
	sources := make([]model.Message, 0, len(messageIDs))
	for _, messageID := range messageIDs {
		if _, ok := seen[messageID]; ok {
			continue
		}
		seen[messageID] = struct{}{}

		source, err := s.messageRepo.GetByID(ctx, messageID)
		if err != nil {
			return nil, err
		}
		if source.DeletedAt != nil || source.ConversationID != sourceConversationID {
			return nil, repository.ErrMessageNotFound
		}
		sources = append(sources, source)
	}

	sort.SliceStable(sources, func(i, j int) bool {
		if sources[i].CreatedAt.Equal(sources[j].CreatedAt) {
			return sources[i].ID.String() < sources[j].ID.String()
		}
		return sources[i].CreatedAt.Before(sources[j].CreatedAt)
	})

	sourceIDs := make([]uuid.UUID, len(sources))
	for i := range sources {
		sourceIDs[i] = sources[i].ID
	}
	atts, err := s.attachmentRepo.ListByMessageIDs(ctx, sourceIDs)
	if err != nil {
		return nil, err
	}
	attsByMessage := map[uuid.UUID][]model.MessageAttachment{}
	for _, att := range atts {
		if att.MessageID == nil {
			continue
		}
		attsByMessage[*att.MessageID] = append(attsByMessage[*att.MessageID], att)
	}

	sender, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	created := make([]model.Message, 0, len(sources))
	for _, source := range sources {
		message, err := model.NewMessage(targetConversationID, source.Content, sender)
		if err != nil {
			return nil, err
		}
		message.ForwardedFromID = &source.ID
		msg, err := s.messageRepo.CreateMessage(ctx, message)
		if err != nil {
			return nil, err
		}
		for _, att := range attsByMessage[source.ID] {
			clone := model.MessageAttachment{
				ID:             uuid.Must(uuid.NewV7()),
				MessageID:      &msg.ID,
				ConversationID: targetConversationID,
				UploaderID:     userID,
				ObjectKey:      att.ObjectKey,
				ThumbKey:       att.ThumbKey,
				Mime:           att.Mime,
				Kind:           att.Kind,
				SizeBytes:      att.SizeBytes,
				Width:          att.Width,
				Height:         att.Height,
				DurationMs:     att.DurationMs,
				OriginalName:   att.OriginalName,
				CreatedAt:      time.Now().UTC(),
			}
			if _, err := s.attachmentRepo.Create(ctx, clone); err != nil {
				return nil, err
			}
		}
		created = append(created, msg)
	}

	if err := s.hydrateMessages(ctx, created); err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) hydrateMessages(ctx context.Context, messages []model.Message) error {
	if len(messages) == 0 {
		return nil
	}
	ids := make([]uuid.UUID, 0, len(messages))
	for i := range messages {
		ids = append(ids, messages[i].ID)
		messages[i].Sender.AvatarURL = s.urls.Resolve(ctx, messages[i].Sender.AvatarURL)
	}
	atts, err := s.attachmentRepo.ListByMessageIDs(ctx, ids)
	if err != nil {
		return err
	}
	byMessage := map[uuid.UUID][]model.MessageAttachment{}
	for _, att := range atts {
		s.resolveAttachment(ctx, &att)
		if att.MessageID == nil {
			continue
		}
		byMessage[*att.MessageID] = append(byMessage[*att.MessageID], att)
	}
	for i := range messages {
		messages[i].Attachments = byMessage[messages[i].ID]
		if messages[i].Attachments == nil {
			messages[i].Attachments = []model.MessageAttachment{}
		}
	}
	return nil
}

func (s *Service) resolveAttachment(ctx context.Context, att *model.MessageAttachment) {
	if att == nil || s.urls == nil {
		return
	}
	att.URL = s.urls.Resolve(ctx, att.ObjectKey)
	if att.ThumbKey != "" {
		att.ThumbURL = s.urls.Resolve(ctx, att.ThumbKey)
	}
}

func classifyMime(mime, filename string) (model.AttachmentKind, string, int64) {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	switch {
	case strings.HasPrefix(mime, "image/"):
		if ext == "" {
			ext = "jpg"
		}
		return model.AttachmentKindImage, ext, maxImageBytes
	case strings.HasPrefix(mime, "video/"):
		if ext == "" {
			ext = "mp4"
		}
		return model.AttachmentKindVideo, ext, maxVideoBytes
	default:
		if ext == "" {
			ext = "bin"
		}
		return model.AttachmentKindFile, ext, maxFileBytes
	}
}

func makeImageThumb(data []byte) (width, height int, thumb []byte, err error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return 0, 0, nil, err
	}
	bounds := img.Bounds()
	width, height = bounds.Dx(), bounds.Dy()
	resized := imaging.Fit(img, 480, 480, imaging.Lanczos)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 75}); err != nil {
		return 0, 0, nil, err
	}
	return width, height, buf.Bytes(), nil
}

// Cursor helpers kept from previous service.go

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
	hydrated := []model.Message{updated}
	_ = s.hydrateMessages(ctx, hydrated)
	return hydrated[0], nil
}

func (s *Service) GenerateCursor(message model.Message) string {
	c := pagination.Cursor{
		ID:        message.ID,
		CreatedAt: message.CreatedAt,
	}
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *Service) GenerateAttachmentCursor(att model.MessageAttachment) string {
	c := pagination.Cursor{
		ID:        att.ID,
		CreatedAt: att.CreatedAt,
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
