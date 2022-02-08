package handlers

import (
	"net/http"

	"github.com/ziemowit141/payment_api/algorithms"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
	"gorm.io/gorm"
)

type Void struct {
	db *gorm.DB
}

func NewVoidHandler(db *gorm.DB) *Void {
	return &Void{db}
}

func (v *Void) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	voidReq := io_structures.NewVoidRequest(r.Body)
	session := algorithms.NewSession(v.db)

	status := session.AuthorizeWithTransactionId(voidReq.Uid)
	if status != "SUCCESS" {
		voidRes := &io_structures.VoidResponse{
			Status:   status,
			Balance:  0.0,
			Currency: "NaN",
		}
		voidRes.ToJSON(rw)
		return
	}

	session.CancelTransaction(voidReq.Uid)
	voidRes := &io_structures.VoidResponse{
		Status:   "SUCCESS",
		Balance:  session.CardNetBalance(),
		Currency: "PLN"}
	voidRes.ToJSON(rw)
}
