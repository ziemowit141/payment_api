package handlers

import (
	"net/http"

	"github.com/ziemowit141/payment_api/algorithms"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
	"gorm.io/gorm"
)

type Refund struct {
	db *gorm.DB
}

func NewRefundHandler(db *gorm.DB) *Refund {
	return &Refund{db}
}

func (c *Refund) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postRefund(rw, r)
	}

	rw.WriteHeader(http.StatusNotImplemented)
}

func (c *Refund) postRefund(rw http.ResponseWriter, r *http.Request) {
	refundReq := io_structures.NewOrderRequest(r.Body)
	session := algorithms.NewSession(c.db)

	status := session.AuthorizeWithTransactionId(refundReq.TransactionId)
	if status != "SUCCESS" {
		voidRes := &io_structures.VoidResponse{
			Status:   status,
			Balance:  0.0,
			Currency: "NaN",
		}
		voidRes.ToJSON(rw)
		return
	}

	status, balance := session.Refund(refundReq.Amount)

	captureRes := &io_structures.OrderResponse{
		Status:   status,
		Balance:  balance,
		Currency: "PLN"}

	captureRes.ToJSON(rw)
}
