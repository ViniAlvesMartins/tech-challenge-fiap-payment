build:
	docker-compose build

run-dev:
	docker-compose up

run-prod:
	docker-compose -f docker-compose.prod.yaml up

mocks:
	docker-compose run payment-dev-app go generate ./...

test:
	docker-compose run payment-dev-app go test ./...

test-coverage:
	docker-compose run payment-dev-app go test -coverprofile cover.out `go list ./... | egrep -v '(/doc|/infra|/src/application/contract/mock|/src/external/handler/http_server)$\'` && go tool cover -html=cover.out

get-coverage:
	docker-compose run -d payment-dev-app go test -coverprofile cover.out `go list ./... | egrep -v '(/doc|/infra|/src/application/contract/mock|/src/external/handler/http_server)$\')` && go tool cover -func cover.out | fgrep total | awk '{print substr($$3, 1, length($$3)-1)}'

run-test:
	$(MAKE) mocks && $(MAKE) test

