package middleware

import (
	"driver/pkg"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFiberMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fiber Middleware Suite")
}

var _ = Describe("Fiber Middleware Suite", Ordered, func() {

	app := fiber.New()
	app.Get("/api/v1/drivers/nearest", IsApiKeyCorrect)

	Context("Middleware", func() {

		When("Api token is invalid", func() {
			It("should be ok", func() {
				var responseObject pkg.ErrorResponse

				// create a request to our endpoint with authorization header
				req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/drivers/nearest?lang=%f&lat=%f&api_key=%s", 10.0, 10.0, "123213"), nil)

				// send the request to the app
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(401))

				// read the response body
				body, err := io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())

				// unmarshal the response body into our struct
				err = json.Unmarshal(body, &responseObject)
				Expect(err).NotTo(HaveOccurred())

				Expect(responseObject.Error).Should(Equal(pkg.ErrInvalidAPIKey.Error()))
			})
		})

		When("Api token is valid", func() {
			It("should be ok", func() {
				// create a request to our endpoint with authorization header
				req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/drivers/nearest?lang=%f&lat=%f&api_key=%s", 10.0, 10.0, pkg.AppConfigs.GetApplicationConfig().API_KEY), nil)

				// send the request to the app
				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).ShouldNot(Equal(401))
			})
		})
	})
})
