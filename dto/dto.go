package dto

import (
	"github.com/google/uuid"
)

// CheckoutMedicationRequest holds information about the checkout request, including a list of items and the total price.
type CheckoutMedicationRequest struct {
	AllItem    []CheckoutItem `json:"allItem"`    // Array of items in the checkout request
	TotalPrice int64          `json:"totalPrice"` // The total price for all items
}

// CheckoutItem represents a single item in a checkout process, including its ID, price, and quantity.
type CheckoutItem struct {
	MedicationID uuid.UUID `json:"medicationId"` // UUID of the medication being purchased
	Name         string    `json:"name"`
	Price        int64     `json:"price"`    // Price of a single unit of the medication
	Quantity     int       `json:"quantity"` // Quantity of the medication being purchased
}
