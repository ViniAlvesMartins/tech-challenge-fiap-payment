build:
	docker-compose build

run-app:
	docker-compose up dev-app-payment

run-worker:
	docker-compose up dev-app-payment-worker

start-infra:
	docker-compose -f docker-compose-infra.yaml up

swagger:
	docker-compose run dev-app-payment swag init -g internal/external/handler/http_server/app.go -o doc/swagger/

mocks:
	docker-compose run dev-app-payment go generate ./...

test:
	docker-compose run dev-app-payment go test ./...

test-coverage:
	docker-compose run dev-app-payment go test -coverprofile cover.out `go list ./... | egrep -v '(/doc|/infra|/src/pkg/uuid/mock|/src/application/contract/mock|/src/external/handler/http_server|/api)$\'` && go tool cover -html=cover.out

get-coverage:
	docker-compose run -d dev-app-payment go test -coverprofile cover.out `go list ./... | egrep -v '(/doc|/infra|/src/pkg/uuid/mock|/src/application/contract/mock|/src/external/handler/http_server|/api)$\')` && go tool cover -func cover.out | fgrep total | awk '{print substr($$3, 1, length($$3)-1)}'

run-test:
	$(MAKE) mocks && $(MAKE) test

