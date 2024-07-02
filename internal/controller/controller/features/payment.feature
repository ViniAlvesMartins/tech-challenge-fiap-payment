Feature: QR Code creation
    In order to pay for an order
    As a Client
    I need to be able to create a payment QR Code

    Scenario: Create a QR code to pay for an order
        When I send a POST request to "/payments"
        Then Status code should be 201
        And A QR code should be returned