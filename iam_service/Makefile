GO_CMD_MAIN = cmd/*.go

update:
	cd ../../.. && make update

lint:
	cd ../../.. && golangci-lint run app/warehouse/warehouse_service/...

run:
	go run $(GO_CMD_MAIN) server

run-kafka:
	go run $(GO_CMD_MAIN) subscribe kafka

run-report:
	go run $(GO_CMD_MAIN) subscribe report-kafka

sqlc:
	sqlc generate

sqlc-lint:
	sqlc vet
	sqlc generate

migrate:
	echo \# make migrate name="$(name)"
	go run $(GO_CMD_MAIN) migrate create $(name)

migrate-up:
	go run $(GO_CMD_MAIN) migrate up

migrate-down-1:
	go run $(GO_CMD_MAIN) migrate down 1

test:
	go test ./...  -count=1 -v -cover

.PHONY = run generate-sqlc migrate migrate-up test

proto:
	cd ../../.. && tekit-protoc -buf=false ./app/warehouse/iam_service/api

mockery-store:
	cd internal/store && mockery --name Repository --case underscore --inpackage --filename repository_mock.go

mockery-core:
	cd internal/wscore && mockery --name WSCore --case underscore --output ../mock/wscore

mockery-adapter:
	cd internal/adapter && mockery --all --case underscore --output ../mock/adapter

mockery-parser:
	cd internal/wscore/parser && mockery --name Parser --case underscore --output ../mock/parser

mockery-producer:
	cd internal/pubsub/producer && mockery --name Producer --case underscore --output ../../mock/producer

mockery-all:
	make mockery-store && \
	 make mockery-core && \
	 make mockery-parser && \
	 make mockery-adapter && \
	 make mockery-producer