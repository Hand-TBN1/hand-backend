package dto

import (
		"github.com/google/uuid"
)

type CheckoutMedicationRequest struct {
	AllItem    []CheckoutItem
	TotalPrice int64
}

type CheckoutItem struct {
	MedicationID uuid.UUID
	Price        int64
	Quantity     int
}