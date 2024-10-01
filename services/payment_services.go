package services

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService struct{}

// CreatePayment handles the creation of a payment request using Snap
func (service *PaymentService) CreatePayment(orderID string, grossAmount int64) (*snap.Response, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: grossAmount,
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "minute",
			Duration: 5,
		},
		Callbacks: &snap.Callbacks{
			Finish: "http://localhost:3000/appointment-history", 
		},
	}

	// Create transaction using the globally set ServerKey and Environment
	resp, err := snap.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
