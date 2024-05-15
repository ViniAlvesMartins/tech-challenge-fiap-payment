build:
	docker-compose build

run-dev:
	docker-compose up

run-prod:
	docker-compose -f docker-compose.prod.yaml up

mocks:
	docker-compose run payment-dev-app go generate ./...

test:
	docker-compose run payment-dev-app go test `go list ./... | grep -v -e mock -e doc -e infra`

test-coverage:
	docker-compose run payment-dev-app go test `go list ./... | grep -v -e mock -e doc -e infra` -coverprofile cover.out  && go tool cover -html=cover.out

run-test:
	$(MAKE) mocks && $(MAKE) test

