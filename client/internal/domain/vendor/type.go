package vendor

import "fbriansyah/client/internal/domain"

type Vendor struct {
	ID       string
	Name     string
	Priority int
	Setting  domain.VendorSetting
	Status   domain.Status
}
