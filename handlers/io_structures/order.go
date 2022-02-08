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

func (a *OrderRequest) FromJSON(r io.Reader) {
	e := json.NewDecoder(r)
	err := e.Decode(a)
	if err != nil {
		panic(err)
	}
}

type OrderResponse struct {
	Status   string  `json:"status"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
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
