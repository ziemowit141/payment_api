package io_structures

import (
	"encoding/json"
	"io"
)

// swagger:model
type VoidRequest struct {
	// Required: true
	Uid string `json:"uid"`
}

// swagger:parameters void
type VoidRequestWrapper struct {
	// in: body
	Body VoidRequest
}

func (a *VoidRequest) FromJSON(r io.Reader) {
	e := json.NewDecoder(r)
	err := e.Decode(a)
	if err != nil {
		panic(err)
	}
}

func NewVoidRequest(body io.ReadCloser) *VoidRequest {
	authReq := new(VoidRequest)
	authReq.FromJSON(body)
	return authReq
}

// swagger:model
type VoidResponse struct {
	Status   string  `json:"status"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
}

func (a *VoidResponse) ToJSON(w io.Writer) {
	e := json.NewEncoder(w)
	err := e.Encode(a)
	if err != nil {
		panic(err)
	}
}

func (a *VoidResponse) FromJSON(data []byte) {
	err := json.Unmarshal(data, a)
	if err != nil {
		panic(err)
	}
}
