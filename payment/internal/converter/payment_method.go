package converter

import (
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

func PaymentMethodToModel(paymentMethod paymentV1.PaymentMethod) model.PaymentMethod {
	switch paymentMethod {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SPB:
		return model.PaymentMethodSPB
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}
