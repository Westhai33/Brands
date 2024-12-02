# start app
PROTO_DIR = ./proto
GEN_DIR = ./
SWAGGER_DIR = ./swagger

# Go-пакеты в проекте
GOPACKAGES := $(shell go list ./...)
.PHONY : run
run:
	go mod download && go run main.go --config="config/dev.yml"

.PHONY : docker-up
docker-up:
	docker-compose -f docker-compose.yml up --build

# Таргет для применения миграций базы данных
migrate:
	@echo "Applying database migrations..."
	goose -dir ./migrations postgres "host=localhost port=5432 user=brands password=pgpwdbrands dbname=brands sslmode=disable" up
	@echo "Database migrations applied successfully."

swagger: generate
	@echo "Generating Swagger specifications..."
	mkdir -p ./internal/api/v1
	protoc -I=$(PROTO_DIR) \
		--openapi_out=./internal/api/v1/ \
		$(PROTO_DIR)/*.proto
	@echo "Swagger specifications generated and saved to swagger.json."

install-proto:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

generate: install-proto
	@echo "Generating Go code from .proto files..."
	@if [ ! -d "$(GEN_DIR)" ]; then mkdir -p "$(GEN_DIR)"; fi
	protoc -I=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go-grpc_out=$(GEN_DIR) \
		--grpc-gateway_out=$(GEN_DIR) \
		--validate_out="lang=go:$(GEN_DIR)" \
		$(PROTO_DIR)/*.proto
	@echo "Proto files have been successfully compiled."