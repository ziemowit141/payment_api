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

var _ = Describe("Void", func() {
	var db *gorm.DB
	var void *handlers.Void
	var request *io_structures.VoidRequest
	var response *io_structures.VoidResponse

	BeforeEach(func() {
		db = database.NewTestDb()
		void = handlers.NewVoidHandler(db)
		request = &io_structures.VoidRequest{
			Uid: "35c64636-70bc-42b5-593e-d051e29c02bc",
		}
	})

	AfterEach(func() {
		database.DropSchema(db)
	})

	makeRequest := func() {
		body_bytes, err := json.Marshal(request)
		Expect(err).To(BeNil())

		requestBody := string(body_bytes)
		req := httptest.NewRequest(http.MethodPost, "/void", strings.NewReader(requestBody))

		requestRecorder := httptest.NewRecorder()
		void.ServeHTTP(requestRecorder, req)

		res := requestRecorder.Result()
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		response = &io_structures.VoidResponse{}
		response.FromJSON(data)
	}

	Context("with correct data", Ordered, func() {
		BeforeAll(func() {
			makeRequest()
		})
		It("should return 'SUCCESS' status", func() {
			Expect(response.Status).To(Equal("SUCCESS"))
		})
		It("should return account balance", func() {
			Expect(response.Balance).To(Equal(float32(10000.0)))
		})
		It("should return PLN currency", func() {
			Expect(response.Currency).To(Equal("PLN"))
		})
	})

	Context("with non existing id", Ordered, func() {
		BeforeAll(func() {
			request.Uid = "non existing"
			makeRequest()
		})
		It("should return 'WRONG TRANSACTION ID' status", func() {
			Expect(response.Status).To(Equal("WRONG TRANSACTION ID"))
		})
		It("should return zero account balance", func() {
			Expect(response.Balance).To(BeZero())
		})
		It("should return NaN currency", func() {
			Expect(response.Currency).To(Equal("NaN"))
		})
	})
})
