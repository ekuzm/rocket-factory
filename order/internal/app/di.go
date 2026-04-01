package app

import (
	"context"
	"fmt"

	api "github.com/ekuzm/rocket-factory/order/internal/api/v1"
	"github.com/ekuzm/rocket-factory/order/internal/config"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/inventory"
	"github.com/ekuzm/rocket-factory/order/internal/integration/grpc/payment"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres"
	"github.com/ekuzm/rocket-factory/order/internal/repository/postgres/transaction"
	"github.com/ekuzm/rocket-factory/order/internal/service"
	"github.com/ekuzm/rocket-factory/platform/pkg/closer"
	orderV1 "github.com/ekuzm/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type di struct {
	api        orderV1.Handler
	service    api.OrderService
	repository service.OrderRepository

	inventoryClient    service.InventoryPort
	paymentClient      service.PaymentPort
	transactionManager service.TransactionManager

	pool *pgxpool.Pool
}

func NewDI() *di {
	return &di{}
}

func (d *di) API(ctx context.Context) (orderV1.Handler, error) {
	if d.api == nil {
		service, err := d.Service(ctx)
		if err != nil {
			return nil, fmt.Errorf("create order service: %w", err)
		}

		d.api = api.New(service)
	}

	return d.api, nil
}

func (d *di) Service(ctx context.Context) (api.OrderService, error) {
	if d.service == nil {
		repo, err := d.Repository(ctx)
		if err != nil {
			return nil, fmt.Errorf("create order repository: %w", err)
		}
		inventoryClient, err := d.InventoryClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("create inventory client: %w", err)
		}
		paymentClient, err := d.PaymentClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("create payment client: %w", err)
		}
		transactionManager, err := d.TransactionManager(ctx)
		if err != nil {
			return nil, fmt.Errorf("transaction manager: %w", err)
		}

		d.service = service.New(repo, inventoryClient, paymentClient, transactionManager)
	}

	return d.service, nil
}

func (d *di) Repository(ctx context.Context) (service.OrderRepository, error) {
	if d.repository == nil {
		pool, err := d.Pool(ctx)
		if err != nil {
			return nil, fmt.Errorf("create db pool: %w", err)
		}

		d.repository = postgres.New(pool)
	}

	return d.repository, nil
}

func (d *di) InventoryClient(ctx context.Context) (service.InventoryPort, error) {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(config.App().GRPC.InventoryAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("create client connection to inventory service: %w", err)
		}

		closer.Add(func(ctx context.Context) error {
			return inventoryConn.Close()
		})

		d.inventoryClient = inventory.New(inventoryV1.NewInventoryServiceClient(inventoryConn))
	}

	return d.inventoryClient, nil
}

func (d *di) PaymentClient(ctx context.Context) (service.PaymentPort, error) {
	if d.paymentClient == nil {
		paymentConn, err := grpc.NewClient(config.App().GRPC.PaymentAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("create client connection to payment service: %w", err)
		}

		closer.Add(func(ctx context.Context) error {
			return paymentConn.Close()
		})

		d.paymentClient = payment.New(paymentV1.NewPaymentServiceClient(paymentConn))
	}

	return d.paymentClient, nil
}

func (d *di) TransactionManager(ctx context.Context) (service.TransactionManager, error) {
	if d.transactionManager == nil {
		pool, err := d.Pool(ctx)
		if err != nil {
			return nil, fmt.Errorf("create transaction manager: %w", err)
		}

		d.transactionManager = transaction.New(pool)
	}

	return d.transactionManager, nil
}

func (d *di) Pool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pool == nil {
		pool, err := pgxpool.New(ctx, config.App().Postgres.URI())
		if err != nil {
			return nil, err
		}

		closer.Add(func(ctx context.Context) error {
			pool.Close()

			return nil
		})

		d.pool = pool
	}

	return d.pool, nil
}
