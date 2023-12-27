package vendor

import "context"

type Repository interface {
	List(ctx context.Context, listVendor []string) ([]Vendor, error)
}
