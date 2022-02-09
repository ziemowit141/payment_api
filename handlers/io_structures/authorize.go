package io_structures

import (
	"encoding/json"
	"io"
)

// swagger:model
type AuthorizationRequest struct {
	// Required: true
	CreditCardNumber string `json:"credit_card_number" validate:"required"`

	// Required: true
	Expiry string `json:"expiry" validate:"required"`

	// Required: true
	CreditCardCVV string `json:"credit_card_cvv" validate:"required"`

	// Required: true
	Amount float32 `json:"amount" validate:"required"`

	// Required: true
	Currency string `json:"base" validate:"required"`
}

// swagger:parameters authorize
type _ struct {
	// in: body
	Body AuthorizationRequest
}

func NewAuthorizationRequest(body io.ReadCloser) *AuthorizationRequest {
	authReq := new(AuthorizationRequest)
	err := authReq.FromJSON(body)
	if err != nil {
		panic(err)
	}

	return authReq
}

func (a *AuthorizationRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

// swagger:model
type AuthorizationResponse struct {
	// Transaction ID, empty if unauthorized
	Uid string `json:"uid"`

	// Status of authorization
	Status string `json:"status"`

	// Account balance, empty if unauthorized
	Balance float32 `json:"balance"`

	// Transaction currency
	Currency string `json:"currency"`
}

func (a *AuthorizationResponse) ToJSON(w io.Writer) {
	e := json.NewEncoder(w)
	err := e.Encode(a)
	if err != nil {
		panic(err)
	}
}

func (a *AuthorizationResponse) FromJSON(data []byte) {
	err := json.Unmarshal(data, a)
	if err != nil {
		panic(err)
	}
}
