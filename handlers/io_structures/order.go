package io_structures

import (
	"encoding/json"
	"io"
)

type OrderRequest struct {
	TransactionId string  `json:"transaction_id"`
	Amount        float32 `json:"amount"`
}

func NewOrderRequest(body io.ReadCloser) *OrderRequest {
	authReq := new(OrderRequest)
	authReq.FromJSON(body)
	return authReq
}

func (a *OrderRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

type OrderResponse struct {
	Status   string  `json:"status"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
}

func (a *OrderResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *OrderResponse) FromJSON(data []byte) {
	json.Unmarshal(data, a)
}
