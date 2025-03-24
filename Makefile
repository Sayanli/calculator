test: ### run test
	go test -v ./... -cover

cover-html: ### run test with coverage and open html report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	del coverage.out
.PHONY: coverage-html

cover: ### run test with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	del coverage.out
.PHONY: coverage

proto-gen: ### generate protobuf files
	protoc -I=protos/proto protos/proto/calculator/calculator.proto --go_out=./protos/gen/go/ --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go/ --go-grpc_opt=paths=source_relative
.PHONY: protogen

swag: ### generate swagger docs
	swag init -g ./internal/controller/httpserver/server.

run: ### run server
	go run ./cmd/main.go

compose-build: ### build docker image
	docker-compose build

compose-up: ### run docker container
	docker-compose up
