package project

import "errors"

var (
	ErrNotFound              = errors.New("project not found")
	ErrAccessDenied          = errors.New("project access denied")
	ErrDocNotFound           = errors.New("project doc not found")
	ErrDecisionNotFound      = errors.New("project decision not found")
	ErrMemberNotFound        = errors.New("user not found")
	ErrNotGroupChat          = errors.New("only group conversations can be linked")
	ErrNotConversationMember = errors.New("not a conversation member")
	ErrOwnerRequired         = errors.New("project owner required")
	ErrMessageNotFound       = errors.New("message not found")
	ErrMessageNotLinked      = errors.New("message conversation is not linked to project")
)
