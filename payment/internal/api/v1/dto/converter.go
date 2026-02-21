package dto

import (
	"github.com/ekuzm/rocket-factory/payment/internal/model"
	paymentV1 "github.com/ekuzm/rocket-factory/shared/pkg/proto/payment/v1"
)

var PaymentMethodToModel = map[paymentV1.PaymentMethod]model.PaymentMethod{
	paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:    model.PaymentMethodUnknown,
	paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:           model.PaymentMethodCard,
	paymentV1.PaymentMethod_PAYMENT_METHOD_SPB:            model.PaymentMethodSPB,
	paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:    model.PaymentMethodCreditCard,
	paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY: model.PaymentMethodInvestorMoney,
}
