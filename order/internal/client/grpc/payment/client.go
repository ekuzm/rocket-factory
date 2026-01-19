package payment

import (
	"google.golang.org/grpc"

	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(conn *grpc.ClientConn) *client {
	return &client{
		generatedClient: paymentV1.NewPaymentServiceClient(conn),
	}
}
