package infobib

import (
	"fbriansyah/client/internal/domain"
	"fbriansyah/client/internal/domain/vendor"
	"fbriansyah/client/pkg/httpbreaker"
)

type InfobibSender struct {
	SmsSender domain.Sender
}

func NewInfobibSender(vendor vendor.Vendor, client httpbreaker.HttpClient) *InfobibSender {
	smsSender := NewSmsSender(vendor, client)

	return &InfobibSender{
		SmsSender: smsSender,
	}
}
