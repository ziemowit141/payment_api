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
	"github.com/ziemowit141/payment_api/handlers"
	"github.com/ziemowit141/payment_api/handlers/io_structures"
)

var _ = Describe("Authorize", func() {
	var db *gorm.DB
	var authorize *handlers.Authorize
	var request *io_structures.AuthorizationRequest
	var response *io_structures.AuthorizationResponse

	BeforeEach(func() {
		db = database.NewTestDb()
		authorize = handlers.NewAuthorizeHandler(db)
		request = &io_structures.AuthorizationRequest{
			CreditCardNumber: "5105105105105103",
			Expiry:           "02/30",
			CreditCardCVV:    "111",
			Amount:           1000,
			Currency:         "PLN",
		}
	})

	AfterEach(func() {
		database.DropSchema(db)
	})

	makeRequest := func() {
		body_bytes, err := json.Marshal(request)
		Expect(err).To(BeNil())

		requestBody := string(body_bytes)
		req := httptest.NewRequest(http.MethodPost, "/authorize", strings.NewReader(requestBody))

		requestRecorder := httptest.NewRecorder()
		authorize.ServeHTTP(requestRecorder, req)

		res := requestRecorder.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		response = &io_structures.AuthorizationResponse{}
		response.FromJSON(data)
	}

	Context("with correct data", Ordered, func() {
		BeforeAll(func() {
			makeRequest()
		})
		It("should return 'SUCCESS' status", func() {
			Expect(response.Status).To(Equal("SUCCESS"))
		})
		It("should transaction id", func() {
			Expect(response.Uid).NotTo(BeEmpty())
		})
		It("should return account balance", func() {
			Expect(response.Balance).To(Equal(float32(9000)))
		})
		It("should return PLN currency", func() {
			Expect(response.Currency).To(Equal("PLN"))
		})
	})

	Context("expired date", Ordered, func() {
		BeforeAll(func() {
			request.Expiry = "01/09"
			makeRequest()
		})
		It("should return empty transaction id", func() {
			Expect(response.Uid).To(BeEmpty())
		})
		It("should return 'CARD EXPIRED' status", func() {
			Expect(response.Status).To(Equal("CARD EXPIRED"))
		})
		It("should return zero account balance", func() {
			Expect(response.Balance).To(BeZero())
		})
		It("should return NaN currency", func() {
			Expect(response.Currency).To(Equal("NaN"))
		})
	})

	Context("with non existing card number", Ordered, func() {
		BeforeAll(func() {
			request.CreditCardNumber = "non existing"
			makeRequest()
		})
		It("should return empty transaction id", func() {
			Expect(response.Uid).To(BeEmpty())
		})
		It("should return 'WRONG CARD NUMBER' status", func() {
			Expect(response.Status).To(Equal("WRONG CARD NUMBER"))
		})
		It("should return zero account balance", func() {
			Expect(response.Balance).To(BeZero())
		})
		It("should return NaN currency", func() {
			Expect(response.Currency).To(Equal("NaN"))
		})
	})

	Context("with wrong cvv", Ordered, func() {
		BeforeAll(func() {
			request.CreditCardCVV = "000"
		})

		It("should return empty transaction id", func() {
			Expect(response.Uid).To(BeEmpty())
		})
		It("should return 'WRONG CARD NUMBER' status", func() {
			Expect(response.Status).To(Equal("WRONG CARD NUMBER"))
		})
		It("should return zero account balance", func() {
			Expect(response.Balance).To(BeZero())
		})
		It("should return NaN currency", func() {
			Expect(response.Currency).To(Equal("NaN"))
		})
	})

	Context("with amount that is too large", Ordered, func() {
		BeforeAll(func() {
			request.Amount = float32(1001)
		})

		It("should return empty transaction id", func() {
			Expect(response.Uid).To(BeEmpty())
		})
		It("should return 'WRONG CARD NUMBER' status", func() {
			Expect(response.Status).To(Equal("WRONG CARD NUMBER"))
		})
		It("should return zero account balance", func() {
			Expect(response.Balance).To(BeZero())
		})
		It("should return NaN currency", func() {
			Expect(response.Currency).To(Equal("NaN"))
		})
	})
})
