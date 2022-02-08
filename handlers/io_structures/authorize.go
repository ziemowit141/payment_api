package io_structures

import (
	"encoding/json"
	"io"
)

type AuthorizationRequest struct {
	CreditCardNumber string  `json:"credit_card_number"`
	CreditCardCVV    string  `json:"credit_card_cvv"`
	Amount           float32 `json:"amount"`
	Currency         string  `json:"base_currency"`
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

type AuthorizationResponse struct {
	Uid      string  `json:"uid"`
	Status   string  `json:"status"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
}

func (a *AuthorizationResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *AuthorizationResponse) FromJSON(data []byte) {
	err := json.Unmarshal(data, a)
	if err != nil {
		panic(err)
	}
}
