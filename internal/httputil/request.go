package httputil

import (
	"net/http"

	"github.com/google/uuid"
)

func ParseConversationID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("conversationID"))
}
