package arangodb

import (
	"context"
	"errors"
	"fbriansyah/client/internal/domain/senderlog"
	"time"

	"github.com/arangodb/go-driver"
)

var (
	ErrCreateSenderlog = errors.New("failed to create senderlog")
)

type SenderlogRepository struct {
	db         driver.Database
	collection driver.Collection
	columnName string
}

// Create implements senderlog.Repository.
func (r *SenderlogRepository) Create(ctx context.Context, log *senderlog.Senderlog) error {
	var model CreateSenderLogModel
	model.FromDomain(*log)

	_, err := r.collection.CreateDocument(ctx, model)
	if err != nil {
		return err
	}

	return nil
}

func NewSenderlogRepository(db driver.Database, columnName string) (*SenderlogRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	isCollectionExist, err := db.CollectionExists(ctx, columnName)
	if err != nil {
		return nil, err
	}

	if !isCollectionExist {
		setupSenderLogCollection(ctx, db, columnName)
	}

	coll, err := db.Collection(ctx, columnName)
	if err != nil {
		return nil, err
	}

	return &SenderlogRepository{
		db:         db,
		collection: coll,
		columnName: columnName,
	}, nil
}

var _ senderlog.Repository = (*SenderlogRepository)(nil)
