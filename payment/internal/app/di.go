package app

import (
	"context"
	"fmt"

	api "github.com/ekuzm/rocket-factory/payment/internal/api/v1"
	"github.com/ekuzm/rocket-factory/payment/internal/service"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type di struct {
	api     paymentV1.PaymentServiceServer
	service api.PaymentService
}

func NewDI() *di {
	return &di{}
}

func (d *di) API(ctx context.Context) (paymentV1.PaymentServiceServer, error) {
	if d.api == nil {
		service, err := d.Service(ctx)
		if err != nil {
			return nil, fmt.Errorf("create payment service: %w", err)
		}

		d.api = api.New(service)
	}

	return d.api, nil
}

func (d *di) Service(_ context.Context) (api.PaymentService, error) {
	if d.service == nil {
		d.service = service.New()
	}

	return d.service, nil
}
