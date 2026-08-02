package message

import (
	"megan-messenger/internal/model"
)

type GetPagingResponse struct {
	Messages   []model.Message `json:"messages"`
	NextCursor string          `json:"nextCursor"`
}

type SendMessageRequest struct {
	Content string `json:"content" validate:"required,min=1,max=5000"`
}

type SendMessageResponse struct {
	Message model.Message `json:"message"`
}
