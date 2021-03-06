package handlers

import (
	"net/http"

	"github.com/ziemowit141/payment_api/algorithms"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
	"gorm.io/gorm"
)

type Capture struct {
	db *gorm.DB
}

func NewCaptureHandler(db *gorm.DB) *Capture {
	return &Capture{db}
}

func (c *Capture) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		c.postCapture(rw, r)
	}

	rw.WriteHeader(http.StatusNotImplemented)
}

// swagger:route POST /capture payment_api capture
//
// Evaluates eligibility to capture given
// amount from customers account and performs
// this capture if possible
//
//	   Consumes:
//	   - application/json
//
//	   Produces:
//	   - application/json
//
//	   Schemes: http
//
//	   Deprecated: false
//
// 	   Responses:
//       200: OrderResponse
//	     400: OrderResponse
// 		 401: OrderResponse
//       501: description:NotImplemented
func (c *Capture) postCapture(rw http.ResponseWriter, r *http.Request) {
	captureReq := io_structures.NewOrderRequest(r.Body)
	session := algorithms.NewSession(c.db)

	status := session.AuthorizeWithTransactionId(captureReq.TransactionId)
	if status != "SUCCESS" {
		voidRes := &io_structures.VoidResponse{
			Status:   status,
			Balance:  0.0,
			Currency: "NaN",
		}
		rw.WriteHeader(http.StatusUnauthorized)
		voidRes.ToJSON(rw)
		return
	}

	status, balance := session.Capture(captureReq.Amount)

	captureRes := &io_structures.OrderResponse{
		Status:   status,
		Balance:  balance,
		Currency: "PLN"}

	if status != "SUCCESS" {
		rw.WriteHeader(http.StatusBadRequest)
	}

	captureRes.ToJSON(rw)
}
