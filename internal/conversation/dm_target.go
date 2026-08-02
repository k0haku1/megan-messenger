package conversation

import (
	"megan-messenger/internal/httputil"

	"github.com/google/uuid"
)

type dmTarget struct {
	userID   *uuid.UUID
	username *string
}

func parseDMTarget(userID *uuid.UUID, username *string) (target dmTarget, fieldErrors httputil.FieldErrors, ok bool) {
	hasUserID := userID != nil && *userID != uuid.Nil
	hasUsername := username != nil && *username != ""
	if hasUserID == hasUsername {
		return dmTarget{}, httputil.FieldErrors{
			"userId":   "Provide either userId or username",
			"username": "Provide either userId or username",
		}, false
	}

	return dmTarget{userID: userID, username: username}, nil, true
}
