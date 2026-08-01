package notification

import (
	"context"
	"fmt"
	"log/slog"
)

type SmsSender interface {
	SendCode(ctx context.Context, phone, code string) error
}

type LogSmsSender struct {
	logger *slog.Logger
}

func NewLogSmsSender(logger *slog.Logger) *LogSmsSender {
	return &LogSmsSender{logger: logger}
}

func (s *LogSmsSender) SendCode(ctx context.Context, phone, code string) error {
	s.logger.Info("sent SMS verification code", "phone", phone, "code", code)
	fmt.Printf("==================================================\n")
	fmt.Printf("SMS to: %s\n", phone)
	fmt.Printf("Your verification code is: %s\n", code)
	fmt.Printf("==================================================\n")
	return nil
}
