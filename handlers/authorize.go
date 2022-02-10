// Package classification of Payment API
//
// Documentation for payment api
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
package handlers

import (
	"net/http"

	"github.com/ziemowit141/payment_api/algorithms"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
	"gorm.io/gorm"
)

type Authorize struct {
	db *gorm.DB
}

func NewAuthorizeHandler(db *gorm.DB) *Authorize {
	return &Authorize{db}
}

func (a *Authorize) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		a.postAuthorize(rw, r)
	}

	rw.WriteHeader(http.StatusNotImplemented)
}

// swagger:route POST /authorize payment_api authorize
//
// Validates passed data and returns
// transactionId for further transactions
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
//       200: AuthorizationResponse
//       400: AuthorizationResponse
//       401: AuthorizationResponse
//       501: description:NotImplemented
func (a *Authorize) postAuthorize(rw http.ResponseWriter, r *http.Request) {
	authReq := io_structures.NewAuthorizationRequest(r.Body)
	session := algorithms.NewSession(a.db)

	status := session.AuthorizeWithCardDetails(authReq)
	if status != "SUCCESS" {
		authRes := &io_structures.AuthorizationResponse{
			Uid:      "",
			Status:   status,
			Balance:  0.0,
			Currency: "NaN"}
		if status == "UNSUFFICIENT FUNDS" {
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
		}
		authRes.ToJSON(rw)
		return
	}

	transactionId := session.AddTransaction(authReq.Amount)

	authRes := &io_structures.AuthorizationResponse{
		Uid:      transactionId,
		Status:   status,
		Balance:  session.CardNetBalance(),
		Currency: authReq.Currency}
	authRes.ToJSON(rw)
}
