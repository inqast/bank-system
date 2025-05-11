# Название исполняемого файла
BINARY_NAME=./bin/app

# Путь к main.go
MAIN_PACKAGE=./cmd

LOCAL_BIN=$(CURDIR)/bin

# Сборка приложения
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Запуск приложения
start:
	$(BINARY_NAME)

run: build start

up:
	docker compose up -d

down:
	docker compose down

down-v:
	docker compose down -v

# Генерация gRPC и других файлов
.PHONY: generate
generate: 
	mkdir -p pkg/api
	protoc -I api \
		-I vendor.protogen \
		-I vendor.protogen/google/api \
		-I vendor.protogen/validate \
		api/api.proto \
		--go_out=./pkg/api \
		--go_opt=paths=source_relative \
		--go-grpc_out=./pkg/api \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out ./pkg/api \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out=./pkg/api \
		--validate_out="lang=go,paths=source_relative:./pkg/api"

# Скачивание и установка нужных proto-файлов
.PHONY: proto-deps
proto-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
