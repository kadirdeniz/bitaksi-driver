package handler

import (
	"bytes"
	"driver/internal"
	"driver/pkg"
	"driver/test/mock"
	"driver/tools/fiber/middleware"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/golang/mock/gomock"
	"io"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var testingObj *testing.T

func TestFiber(t *testing.T) {
	testingObj = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fiber Suite")
}

var _ = Describe("Fiber Suite", func() {
	ctrl := gomock.NewController(testingObj)
	defer ctrl.Finish()

	var mockRepository *mock.MockIRepository

	BeforeEach(func() {
		mockRepository = mock.NewMockIRepository(ctrl)
	})

	Context("Create a new handler", func() {
		It("Should return a new handler", func() {
			handler := NewHandler(mockRepository)
			Expect(handler).NotTo(BeNil())
		})
	})

	Context("Get nearest drivers", func() {
		When("Drivers are found", func() {
			It("Should return a list of drivers", func() {
				var responseObject internal.NearestLocationResponse

				mockRepository.EXPECT().FindNearestDriver(gomock.Any()).Return(mock.LocationModel, nil).Times(1)

				handler := NewHandler(mockRepository)

				app := fiber.New()
				app.Use(recover.New())
				app.Get("/api/v1/drivers/nearest", middleware.IsApiKeyCorrect, handler.GetNearestDriver)

				// create a request to our endpoint
				req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/drivers/nearest?long=%f&lat=%f&api_key=%s", 10.0, 10.0, pkg.AppConfigs.GetApplicationConfig().API_KEY), nil)
				req.Header.Set("Content-Type", "application/json")

				// send the request to the app
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(200))

				// read the response body
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				// unmarshal the response body into our struct
				err = json.Unmarshal(body, &responseObject)
				Expect(err).NotTo(HaveOccurred())

				Expect(responseObject).Should(Equal(internal.NearestLocationResponse{
					ID:          mock.LocationModel.ID,
					Distance:    mock.LocationModel.Distance,
					Coordinates: mock.LocationModel.Location.Coordinates,
				}))
			})
		})

		When("Drivers are not found", func() {
			It("Should return an error", func() {
				mockRepository.EXPECT().FindNearestDriver(gomock.Any()).Return(internal.Model{}, pkg.ErrDriverNotFound).Times(1)

				handler := NewHandler(mockRepository)

				app := fiber.New()
				app.Use(recover.New())
				app.Get("/api/v1/drivers/nearest", handler.GetNearestDriver)

				// create a request to our endpoint with authorization header
				req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/drivers/nearest?long=%f&lat=%f", 10.0, 10.0), nil)
				req.Header.Set("Content-Type", "application/json")
				// send the request to the app
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(404))

				// read the response body
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				// unmarshal the response body into our struct
				var responseObject pkg.ErrorResponse
				err = json.Unmarshal(body, &responseObject)
				Expect(err).NotTo(HaveOccurred())
				Expect(responseObject.Error).Should(Equal(pkg.ErrDriverNotFound.Error()))
			})
		})
	})

	Context("Bulk create drivers", func() {
		When("Driver request more than 100", func() {
			It("Should return an error", func() {
				var responseObject pkg.ErrorResponse

				handler := NewHandler(mockRepository)

				app := fiber.New()
				app.Use(recover.New())
				app.Post("/api/v1/drivers", middleware.SetContentType, handler.BulkCreateDrivers)

				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(mock.BulkCreateDriversRequest(101))

				// create a request to our endpoint
				req := httptest.NewRequest("POST", "/api/v1/drivers", &buf)
				req.Header.Set("Content-Type", "application/json")

				// send the request to the app
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(400))

				// read the response body
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				// unmarshal the response body into our struct
				err = json.Unmarshal(body, &responseObject)
				Expect(err).NotTo(HaveOccurred())

				Expect(responseObject.Error).Should(Equal(pkg.ErrBulkCreateDriversLimit.Error()))
			})
		})

		When("Driver request less than 100", func() {
			It("Should return a success response", func() {
				mockRepository.EXPECT().BulkCreateDrivers(gomock.Any()).Return(nil).Times(1)

				handler := NewHandler(mockRepository)

				app := fiber.New()
				app.Use(recover.New())
				app.Post("/api/v1/drivers", middleware.SetContentType, handler.BulkCreateDrivers)

				var buf bytes.Buffer
				err := json.NewEncoder(&buf).Encode(mock.BulkCreateDriversRequest(99))

				// create a request to our endpoint
				req := httptest.NewRequest("POST", "/api/v1/drivers", &buf)
				req.Header.Set("Content-Type", "application/json")

				// send the request to the app
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(200))

				// read the response body
				_, err = io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Context("ErrorHandler", func() {
		var responseObject pkg.ErrorResponse

		It("Should return an error", func() {
			handler := NewHandler(mockRepository)
			Expect(handler).NotTo(BeNil())

			app := fiber.New(fiber.Config{
				ErrorHandler: handler.ErrorHandler,
			})
			app.Use(recover.New())

			app.Get("/api/v1/error", func(c *fiber.Ctx) error {
				panic("error")
			})

			req := httptest.NewRequest("GET", "/api/v1/error", nil)
			req.Header.Set("X-Request-Id", "Test-Header")
			resp, err := app.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(500))

			responseBody, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())

			err = json.Unmarshal(responseBody, &responseObject)
			Expect(err).To(BeNil())

			Expect(responseObject.Error).To(Equal(pkg.ErrInternalServer.Error()))
		})
	})
})
