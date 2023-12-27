package arangodb

import (
	"fbriansyah/client/internal/domain"
	"fbriansyah/client/internal/domain/vendor"
)

type VendorModel struct {
	ID       string               `json:"_id"`
	Key      string               `json:"_key"`
	Name     string               `json:"name"`
	Priority int                  `json:"priority"`
	Setting  domain.VendorSetting `json:"setting"`
	Status   int                  `json:"status"`
}

type CreateVendorModel struct {
	Key      string               `json:"_key"`
	Name     string               `json:"name"`
	Priority int                  `json:"priority"`
	Setting  domain.VendorSetting `json:"setting"`
	Status   int                  `json:"status"`
}

func (m *VendorModel) ToDomain() vendor.Vendor {
	return vendor.Vendor{
		ID:       m.Key,
		Name:     m.Name,
		Priority: m.Priority,
		Setting:  m.Setting,
		Status:   domain.Status(m.Status),
	}
}
