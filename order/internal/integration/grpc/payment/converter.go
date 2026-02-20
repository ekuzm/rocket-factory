package payment

import (
	"github.com/ekuzm/rocket-factory/order/internal/model"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

var paymentMethodToGRPC = map[model.PaymentMethod]paymentV1.PaymentMethod{
	model.PaymentMethodUnknown:       paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
	model.PaymentMethodCard:          paymentV1.PaymentMethod_PAYMENT_METHOD_CARD,
	model.PaymentMethodSPB:           paymentV1.PaymentMethod_PAYMENT_METHOD_SPB,
	model.PaymentMethodCreditCard:    paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	model.PaymentMethodInvestorMoney: paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
}
