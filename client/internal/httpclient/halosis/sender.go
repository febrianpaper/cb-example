package halosis

import (
	"fbriansyah/client/internal/domain"
	"fbriansyah/client/internal/domain/vendor"
	"fbriansyah/client/pkg/httpbreaker"
)

type HalosisSender struct {
	SmsSender domain.Sender
}

func NewHalosisSender(vendor vendor.Vendor, client httpbreaker.HttpClient) *HalosisSender {
	smsSender := NewSmsSender(vendor, client)

	return &HalosisSender{
		SmsSender: smsSender,
	}
}
