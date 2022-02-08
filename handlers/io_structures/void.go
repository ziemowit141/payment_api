package io_structures

import (
	"encoding/json"
	"io"
)

type VoidRequest struct {
	Uid string `json:"uid"`
}

func (a *VoidRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

func NewVoidRequest(body io.ReadCloser) *VoidRequest {
	authReq := new(VoidRequest)
	authReq.FromJSON(body)
	return authReq
}

type VoidResponse struct {
	Status   string  `json:"status"`
	Balance  float32 `json:"balance"`
	Currency string  `json:"currency"`
}

func (a *VoidResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *VoidResponse) FromJSON(data []byte) {
	json.Unmarshal(data, a)
}
