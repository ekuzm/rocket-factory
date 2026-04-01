package app

import (
	"context"
	"fmt"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/config"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	mng "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type di struct {
	api        inventoryV1.InventoryServiceServer
	service    api.InventoryService
	repository service.InventoryRepository

	client *mng.Client
	db     *mng.Database
}

func NewDI() *di {
	return &di{}
}

func (d *di) API(ctx context.Context) (inventoryV1.InventoryServiceServer, error) {
	if d.api == nil {
		service, err := d.Service(ctx)
		if err != nil {
			return nil, fmt.Errorf("create inventory service: %w", err)
		}

		d.api = api.New(service)
	}

	return d.api, nil
}

func (d *di) Service(ctx context.Context) (api.InventoryService, error) {
	if d.service == nil {
		repository, err := d.Repository(ctx)
		if err != nil {
			return nil, fmt.Errorf("create inventory repository: %w", err)
		}

		d.service = service.New(repository)
	}

	return d.service, nil
}

func (d *di) Repository(ctx context.Context) (service.InventoryRepository, error) {
	if d.repository == nil {
		database, err := d.Database(ctx)
		if err != nil {
			return nil, fmt.Errorf("create mongo database: %w", err)
		}

		repository, err := d.newRepository(database)
		if err != nil {
			return nil, fmt.Errorf("create mongo repository: %w", err)
		}

		d.repository = repository
	}

	return d.repository, nil
}

func (d *di) Database(ctx context.Context) (*mng.Database, error) {
	if d.db == nil {
		client, err := d.MongoClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("create mongo client: %w", err)
		}

		d.db = client.Database(config.App().Mongo.Name())
	}

	return d.db, nil
}

func (d *di) MongoClient(ctx context.Context) (*mng.Client, error) {
	if d.client == nil {
		connectCtx, connectCancel := context.WithTimeout(ctx, mongoConnectTimeout)
		defer connectCancel()

		client, err := mng.Connect(connectCtx, options.Client().ApplyURI(config.App().Mongo.URI()))
		if err != nil {
			return nil, fmt.Errorf("connect mongo client: %w", err)
		}

		if err = client.Ping(connectCtx, nil); err != nil {
			return nil, fmt.Errorf("ping mongo client: %w", err)
		}

		closer.Add(func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.client = client
	}

	return d.client, nil
}

func (d *di) newRepository(db *mng.Database) (_ service.InventoryRepository, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("init mongo repository: %v", recovered)
		}
	}()

	return mongo.New(db), nil
}
