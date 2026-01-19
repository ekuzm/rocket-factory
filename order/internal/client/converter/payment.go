package converter

import (
	"github.com/ekuzm/rocket-factory/order/internal/model"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(paymentMethod model.PaymentMethod) paymentV1.PaymentMethod {
	switch paymentMethod {
	case model.PaymentMethodCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSPB:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SPB
	case model.PaymentMethodCreditCard:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
