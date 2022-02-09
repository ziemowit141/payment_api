package io_structures

import (
	"encoding/json"
	"io"
)

// Capture or Order request for transaction
// swagger:model
type OrderRequest struct {
	// Required: true
	TransactionId string `json:"transaction_id"`

	// Required: true
	Amount float32 `json:"amount"`
}

// swagger:parameters capture refund
type _ struct {
	// in: body
	Body OrderRequest
}

func NewOrderRequest(body io.ReadCloser) *OrderRequest {
	authReq := new(OrderRequest)
	authReq.FromJSON(body)
	return authReq
}

func (a *OrderRequest) FromJSON(r io.Reader) {
	e := json.NewDecoder(r)
	err := e.Decode(a)
	if err != nil {
		panic(err)
	}
}

// swagger:model
type OrderResponse struct {
	// Status of Order Request (Capture or Refund)
	Status string `json:"status"`

	// Current Account balance
	Balance float32 `json:"balance"`

	// Currency at account
	Currency string `json:"currency"`
}

func (a *OrderResponse) ToJSON(w io.Writer) {
	e := json.NewEncoder(w)
	err := e.Encode(a)
	if err != nil {
		panic(err)
	}
}

func (a *OrderResponse) FromJSON(data []byte) {
	err := json.Unmarshal(data, a)
	if err != nil {
		panic(err)
	}
}
