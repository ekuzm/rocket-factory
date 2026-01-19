package v1

import (
	"github.com/ekuzm/rocket-factory/payment/internal/service"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewAPI(service service.PaymentService) *api {
	return &api{
		service: service,
	}
}
