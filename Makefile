# ╔════════════════════════════════════════════════════════════════════════╗
#                          ⭐ Start App ⭐
# ╚════════════════════════════════════════════════════════════════════════╝
.PHONY: run
run:
	go mod download && go run main.go --config="config/dev.yml"

.PHONY: run-local
run-local:
	go mod download && go run main.go --config="config/local.yml"

.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yml up --build

# ╔════════════════════════════════════════════════════════════════════════╗
#                 ⚙️ Database Migrations Targets
# ╚════════════════════════════════════════════════════════════════════════╝
# Таргеты для применения миграций базы данных
include .env
DB_STRING := "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable"
# Директория миграций
MIGRATIONS_DIR="./migrations"
## migration-create: Создает новую миграцию
.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATIONS_DIR)" create "$(MIGRATION_NAME)" sql

 ## migration-up: Применяет миграции к базе данных
.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_STRING)" up

 ## migration-down: Откатывает миграции в базе данных
.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATIONS_DIR)" postgres "$(DB_STRING)" down

# ╔════════════════════════════════════════════════════════════════════════╗
#                 ⚙️ Proto Files and Code Generation
# ╚════════════════════════════════════════════════════════════════════════╝
# Директории для .api файлов и сгенерированных файлов
PROTO_DIR = $(CURDIR)/api
GEN_DIR:=$(CURDIR)/internal/pb
# Путь к локальным бинарным файлам Go
LOCAL_BIN:=$(CURDIR)/bin
# Список всех .api файлов в директории PROTO_DIR и ее поддиректориях
PROTO_FILES=$(shell find $(PROTO_DIR) -name "*.proto")


swagger: generate
	@echo "Generating Swagger specifications..."
	mkdir -p ./internal/pb
	protoc -I=$(PROTO_DIR) \
		--openapi_out=./internal/pb/ \
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