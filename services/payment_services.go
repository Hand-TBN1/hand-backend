package services

import (
	"github.com/veritrans/go-midtrans"
)

type PaymentService struct {
	MidtransClient midtrans.Client
}

func NewPaymentService(client midtrans.Client) *PaymentService {
	return &PaymentService{
		MidtransClient: client,
	}
}

func (service *PaymentService) CreatePayment(orderID string, grossAmount int64) (*midtrans.SnapResponse, error) {
	snapGateway := midtrans.SnapGateway{
		Client: service.MidtransClient,
	}

	req := &midtrans.SnapReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
	}

	resp, err := snapGateway.GetToken(req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
