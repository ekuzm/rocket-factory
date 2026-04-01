package app

import (
	"context"
	"fmt"

	mng "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	api "github.com/ekuzm/rocket-factory/inventory/internal/api/v1"
	"github.com/ekuzm/rocket-factory/inventory/internal/config"
	"github.com/ekuzm/rocket-factory/inventory/internal/repository/mongo"
	"github.com/ekuzm/rocket-factory/inventory/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
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

func (d *di) Repository(ctx context.Context) (_ service.InventoryRepository, err error) {
	if d.repository == nil {
		db, err := d.Database(ctx)
		if err != nil {
			return nil, fmt.Errorf("create mongo database: %w", err)
		}

		d.repository = mongo.New(ctx, db)
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
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
