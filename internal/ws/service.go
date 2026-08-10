package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"megan-messenger/internal/model"
	"megan-messenger/internal/repository"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	rdb              *redis.Client
	upgrader         *websocket.Upgrader
	userRepo         repository.UserRepository
	messageRepo      repository.MessageRepository
	conversationRepo repository.ConversationRepository
}

func NewService(
	rdb *redis.Client,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
	conversationRepo repository.ConversationRepository,
	allowedOrigins []string,
) *Service {
	return &Service{
		rdb:              rdb,
		userRepo:         userRepo,
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")

				for _, allowedOrigin := range allowedOrigins {
					if origin == allowedOrigin {
						return true
					}
				}
				return false
			},
		},
	}
}

func (s *Service) HandleWebSocket(
	w http.ResponseWriter,
	r *http.Request,
	conversation model.Conversation,
	userID uuid.UUID,
) error {
	user, err := s.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		return err
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := s.rdb.Subscribe(ctx, conversation.ID.String())
	defer sub.Close()

	inErr := make(chan error, 1)
	outErr := make(chan error, 1)

	go s.handleIncoming(ctx, conn, conversation.ID, user, inErr)
	go s.handleOutgoing(ctx, conn, sub.Channel(), outErr)

	select {
	case err := <-inErr:
		return err
	case err := <-outErr:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// HandleInboxWebSocket streams conversation activity for the chat list.
func (s *Service) HandleInboxWebSocket(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := s.rdb.Subscribe(ctx, UserChannel(userID))
	defer sub.Close()

	outErr := make(chan error, 1)
	go s.handleOutgoing(ctx, conn, sub.Channel(), outErr)
	go s.drainIncoming(ctx, conn, outErr)

	select {
	case err := <-outErr:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Service) drainIncoming(ctx context.Context, conn *websocket.Conn, errChan chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, _, err := conn.ReadMessage(); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					errChan <- fmt.Errorf("inbox websocket closed unexpectedly: %w", err)
				}
				return
			}
		}
	}
}

func (s *Service) handleIncoming(
	ctx context.Context,
	conn *websocket.Conn,
	conversationID uuid.UUID,
	user model.User,
	errChan chan error,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgType, msgBytes, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					errChan <- fmt.Errorf("websocket closed unexpectedly: %w", err)
				}
				return
			}
			if msgType != websocket.TextMessage {
				continue
			}

			message, err := s.processMessage(ctx, conversationID, string(msgBytes), user)
			if err != nil {
				slog.Warn("Error creating message", "err", err)
				continue
			}

			if err := s.PublishMessage(ctx, conversationID, message); err != nil {
				slog.Warn("Error publishing message", "err", err)
			}
			if err := s.PublishInboxActivity(ctx, message); err != nil {
				slog.Warn("Error publishing inbox activity", "err", err)
			}
		}
	}
}

func (s *Service) processMessage(ctx context.Context, conversationID uuid.UUID, content string, sender model.User) (model.Message, error) {
	msg, err := model.NewMessage(conversationID, content, sender)
	if err != nil {
		return model.Message{}, err
	}

	createdMessage, err := s.messageRepo.CreateMessage(ctx, msg)
	if err != nil {
		return model.Message{}, err
	}

	return createdMessage, nil
}

func (s *Service) PublishMessage(ctx context.Context, conversationID uuid.UUID, message model.Message) error {
	payload, err := json.Marshal(MessageEvent(message))
	if err != nil {
		return err
	}

	return s.rdb.Publish(ctx, conversationID.String(), payload).Err()
}

func (s *Service) PublishRead(ctx context.Context, conversationID, userID uuid.UUID, lastReadAt time.Time) error {
	payload, err := json.Marshal(ReadReceiptEvent(conversationID, userID, lastReadAt))
	if err != nil {
		return err
	}

	return s.rdb.Publish(ctx, conversationID.String(), payload).Err()
}

// PublishInboxActivity notifies each member's inbox channel about a new message.
func (s *Service) PublishInboxActivity(ctx context.Context, message model.Message) error {
	if s.conversationRepo == nil {
		return nil
	}
	memberIDs, err := s.conversationRepo.ListMemberIDs(ctx, message.ConversationID)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(ConversationActivityFromMessage(message))
	if err != nil {
		return err
	}
	for _, memberID := range memberIDs {
		if err := s.rdb.Publish(ctx, UserChannel(memberID), payload).Err(); err != nil {
			slog.Warn("failed to publish inbox activity", "userId", memberID, "err", err)
		}
	}
	return nil
}

func (s *Service) handleOutgoing(
	ctx context.Context,
	conn *websocket.Conn,
	ch <-chan *redis.Message,
	errChan chan error,
) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				errChan <- err
				return
			}
		case msg, ok := <-ch:
			if !ok {
				errChan <- fmt.Errorf("redis channel closed")
				return
			}
			err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			if err != nil {
				errChan <- err
				return
			}
		}
	}
}
