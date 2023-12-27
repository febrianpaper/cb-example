package arangodb

import (
	"context"
	"fbriansyah/client/internal/domain"

	"github.com/arangodb/go-driver"
)

// Setup new vendor collection and insert default data
func setupVendorCollection(ctx context.Context, db driver.Database, columnName string) error {
	collection, err := db.CreateCollection(ctx, columnName, nil)
	if err != nil {
		return err
	}
	vendors := []CreateVendorModel{
		{
			Key:      "infobib",
			Name:     "infobib",
			Priority: 1,
			Setting: domain.VendorSetting{
				BaseUrl:           "http://localhost:8081/v1",
				IsSupportTemplate: false,
				SmsEndpoint:       "/sms",
				WAEndpoint:        "/wa",
				AllowSms:          true,
				AllowWa:           true,
			},
			Status: int(domain.StatusActive),
		},
		{
			Key:      "halosis",
			Name:     "halosis",
			Priority: 2,
			Setting: domain.VendorSetting{
				BaseUrl:           "http://localhost:8082/v1",
				IsSupportTemplate: true,
				SmsEndpoint:       "/sms",
				WAEndpoint:        "/wa",
				AllowSms:          true,
				AllowWa:           false,
			},
			Status: int(domain.StatusActive),
		},
	}

	_, _, err = collection.CreateDocuments(ctx, vendors)

	return err
}

func setupSenderLogCollection(ctx context.Context, db driver.Database, columnName string) error {
	_, err := db.CreateCollection(ctx, columnName, nil)
	if err != nil {
		return err
	}
	return nil
}
