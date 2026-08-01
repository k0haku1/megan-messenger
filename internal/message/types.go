package message

import (
	"megan-messenger/internal/model"
)

type GetPagingResponse struct {
	Messages   []model.Message `json:"messages"`
	NextCursor string          `json:"nextCursor"`
}
