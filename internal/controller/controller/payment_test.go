package controller

import (
	"encoding/json"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/config"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/controller/serializer/input"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/entity"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/internal/entities/enum"
	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"log"
	"testing"
)

var requestBody input.PaymentDto
var expectedResponse Response
var responseBody Response
var statusCode int

func iSendAPostRequest() error {
	config, _ := config.NewConfig()
	client := resty.New().SetBaseURL(config.OrdersURL)

	expectedResponse = Response{
		Error: "",
		Data: entity.QRCodePayment{
			QRCode: "00020101021243650016COM.MERCADOLIBRE02013063638f1192a-5fd1-4180-a180-8bcae3556bc35204000053039865802BR5925IZABEL AAAA DE MELO6007BARUERI62070503***63040B6D",
		},
	}

	requestBody = input.PaymentDto{
		Type:    string(enum.PaymentTypeQRCode),
		OrderId: 1,
	}

	httpmock.ActivateNonDefault(client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/payments", httpmock.NewJsonResponderOrPanic(201, expectedResponse))

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&responseBody).
		SetBody(requestBody).
		Post("/payments")

	if err != nil {
		log.Printf("Request failed: %s", err)
	}

	statusCode = resp.StatusCode()

	return nil
}

func statusCodeShouldBe(expectedStatus int) error {
	if expectedStatus != statusCode {
		if statusCode >= 400 {
			return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, statusCode)
		}
		return fmt.Errorf("expected response code to be: %d, but actual is: %d", expectedStatus, statusCode)
	}

	return nil
}

func aQRCodeShouldBeReturned() error {
	jsonResponse, _ := json.Marshal(responseBody)
	jsonExpectedResponse, _ := json.Marshal(expectedResponse)

	if string(jsonResponse) != string(jsonExpectedResponse) {
		return fmt.Errorf("Expected: %s\nGot: %s", string(jsonExpectedResponse), string(jsonResponse))
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send a POST request to "/payments"$`, iSendAPostRequest)
	ctx.Step(`^Status code should be (\d+)$`, statusCodeShouldBe)
	ctx.Step(`A QR code should be returned$`, aQRCodeShouldBeReturned)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
