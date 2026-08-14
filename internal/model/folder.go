package model

import (
	"time"

	"github.com/google/uuid"
)

type ChatFolder struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon,omitempty"`
	Position  int        `json:"position"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func NewChatFolder(userID uuid.UUID, name, icon string, position int) ChatFolder {
	now := time.Now()
	return ChatFolder{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    userID,
		Name:      name,
		Icon:      icon,
		Position:  position,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
