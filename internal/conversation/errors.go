package conversation

import "errors"

var ErrInvalidGroupMember = errors.New("group member not found")

var ErrNotGroupConversation = errors.New("not a group conversation")
