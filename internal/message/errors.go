package message

import "errors"

var ErrNotMessageSender = errors.New("only the sender can delete a message for everyone")
