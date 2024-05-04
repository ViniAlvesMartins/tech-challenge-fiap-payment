build:
	docker-compose build

run-dev:
	docker-compose up

run-prod:
	docker-compose -f docker-compose.prod.yaml up

mocks:
	docker-compose run dev-app go generate ./...

test:
	docker-compose run dev-app go test -v ./...

run-test:
	$(MAKE) mocks && $(MAKE) test
