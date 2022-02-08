package handlers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/ziemowit141/payment_api/database"
	"github.com/ziemowit141/payment_api/database/models"
	"github.com/ziemowit141/payment_api/handlers"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
)

var _ = Describe("Refund", func() {
	var db *gorm.DB
	var refund *handlers.Refund
	var request *io_structures.OrderRequest
	var response *io_structures.OrderResponse

	BeforeEach(func() {
		db = database.NewTestDb()
		refund = handlers.NewRefundHandler(db)
		request = &io_structures.OrderRequest{
			TransactionId: "35c64636-70bc-42b5-593e-d051e29c02bc",
			Amount:        1000,
		}
	})

	AfterEach(func() {
		database.DropSchema(db)
	})

	makeRequest := func() {
		body_bytes, err := json.Marshal(request)
		Expect(err).To(BeNil())

		requestBody := string(body_bytes)
		req := httptest.NewRequest(http.MethodPost, "/refund", strings.NewReader(requestBody))

		requestRecorder := httptest.NewRecorder()
		refund.ServeHTTP(requestRecorder, req)

		res := requestRecorder.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		response = &io_structures.OrderResponse{}
		response.FromJSON(data)
	}

	Context("with correct data", Ordered, func() {
		BeforeAll(func() {
			capture := &models.Capture{
				TransactionId: "35c64636-70bc-42b5-593e-d051e29c02bc",
				Amount:        1000.0,
			}
			result := db.Create(capture)
			Expect(result.Error).To(BeNil())
			makeRequest()
		})
		It("should return 'SUCCESS' status", func() {
			Expect(response.Status).To(Equal("SUCCESS"))
		})
		It("should return account balance", func() {
			// It is 1000 more because we created capture as a record instead of through API
			Expect(response.Balance).To(Equal(float32(6000)))
		})
		It("should return PLN currency", func() {
			Expect(response.Currency).To(Equal("PLN"))
		})
	})

	Context("with amount larger than order value", Ordered, func() {
		BeforeAll(func() {
			largerThanOrderValue := float32(100000.0)
			request.Amount = largerThanOrderValue
			makeRequest()
		})
		It("should return 'REFUND CANT BE LARGER THAN CAPTURED VALUE status", func() {
			Expect(response.Status).To(Equal("REFUND CANT BE LARGER THAN CAPTURED VALUE"))
		})
		It("should return actual account balance", func() {
			Expect(response.Balance).To(Equal(float32(5000.0)))
		})
		It("should return PLN currency", func() {
			Expect(response.Currency).To(Equal("PLN"))
		})
	})

	Context("with wrong transaction id", Ordered, func() {
		BeforeAll(func() {
			request.TransactionId = "wrong"
			makeRequest()
		})
		It("should return 'WRONG TRANSACTION ID' status", func() {
			Expect(response.Status).To(Equal("WRONG TRANSACTION ID"))
		})
		It("should return zero account balance", func() {
			Expect(response.Balance).To(BeZero())
		})
		It("should return PLN currency", func() {
			Expect(response.Currency).To(Equal("NaN"))
		})
	})
})
