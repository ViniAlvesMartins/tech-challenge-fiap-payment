package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DynamoDBRegion       string `envconfig:"dynamodb_region"`
	DynamoDBUrl          string `envconfig:"dynamodb_url"`
	DynamoDBAccessKey    string `envconfig:"dynamodb_access_key"`
	DynamoDBSecretAccess string `envconfig:"dynamodb_secret_access"`

	SnsRegion                string `envconfig:"sns_region"`
	SnsUrl                   string `envconfig:"sns_url"`
	SnsAccessKey             string `envconfig:"sns_access_key"`
	UpdatePaymentStatusTopic string `envconfig:"payment_status_updated_topic"`
	ProductionFailedQueue    string `envconfig:"production_failed_queue"`
	SnsSecretAccess          string `envconfig:"sns_secret_access"`

	OrdersURL string `envconfig:"orders_url"`
}

func NewConfig() (cfg Config, err error) {
	err = envconfig.Process("", &cfg)
	return
}
