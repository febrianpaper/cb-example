package arangodb

import (
	"context"
	"fbriansyah/client/internal/domain/vendor"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
)

type (
	VendorRepository struct {
		db         driver.Database
		collection driver.Collection
		columnName string
	}
)

// List Get All Vendor
func (r *VendorRepository) List(ctx context.Context, listVendor []string) ([]vendor.Vendor, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	query := fmt.Sprintf(`FOR v IN %s filter v._key in @keys SORT v.priority RETURN v`, r.columnName)

	cursor, err := r.db.Query(ctx, query, map[string]any{"keys": listVendor})
	if err != nil {
		return []vendor.Vendor{}, err
	}
	defer cursor.Close()

	var vendors []vendor.Vendor

	for {
		var vm VendorModel
		_, err := cursor.ReadDocument(ctx, &vm)
		if driver.IsNoMoreDocuments(err) {
			break
		}
		if err != nil {
			return vendors, err
		}
		vendors = append(vendors, vm.ToDomain())
	}

	return vendors, nil
}

func NewVendorRepository(db driver.Database, columnName string) (*VendorRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	isCollectionExist, err := db.CollectionExists(ctx, columnName)
	if err != nil {
		return nil, err
	}

	if !isCollectionExist {
		err := setupVendorCollection(ctx, db, columnName)
		if err != nil {
			return nil, err
		}
	}

	coll, err := db.Collection(ctx, columnName)
	if err != nil {
		return nil, err
	}

	return &VendorRepository{
		db:         db,
		collection: coll,
		columnName: columnName,
	}, nil
}

var _ vendor.Repository = (*VendorRepository)(nil)
