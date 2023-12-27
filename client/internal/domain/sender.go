package domain

import (
	"context"
	"fbriansyah/client/internal/domain/senderlog"
)

type TemplateConfig struct {
	Name string
	Data []string
}

type Sender interface {
	// Send sending sms with raw message.
	Send(ctx context.Context, to string, message string) (*senderlog.Senderlog, error)
	// SendWithTemplate sending sms with template
	SendWithTemplate(ctx context.Context, to string, cfg TemplateConfig) (*senderlog.Senderlog, error)
}
