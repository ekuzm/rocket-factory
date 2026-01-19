package payment

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/ekuzm/rocket-factory/payment/internal/dto"
	"github.com/ekuzm/rocket-factory/payment/internal/model"
)

func (s *service) PayOrder(ctx context.Context, input dto.PayOrderInput) (dto.PayOrderOutput, error) {
	var output dto.PayOrderOutput

	_, err := model.NewPayment(input.OrderUUID, input.UserUUID, input.PaymentMethod)
	if err != nil {
		return output, fmt.Errorf("payment validation: %w", err)
	}

	transactionUUID := uuid.New().String()

	log.Printf("Payment was successfully, transaction uuid: %s", transactionUUID)

	return dto.PayOrderOutput{
		TransactionUUID: transactionUUID,
	}, nil
}
