package mongodb

import (
	"context"
	"fmt"
	"github.com/ViniAlvesMartins/tech-challenge-fiap-payment/infra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewConnection(cfg infra.Config) (*mongo.Database, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	fmt.Println("conectou com sucesso")

	return client.Database("ze_burguer_payment"), nil
}
