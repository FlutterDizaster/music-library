BINARY_NAME=mlibrary
BINARY_PATH=./build

CMD_DIR=./cmd

SWAGGER_OUTPUT=./swagger
SWAGGER_DIRS=./internal/server/handler,./internal/models
SWAGGER_GENERAL=handler.go

.PHONY: all
all: build

# Tests
.PHONY: test
test:
	go test ./... -v

# Code generation
.PHONY: generate
generate:
	go generate ./...

# Building
.PHONY: build
build: generate
	go build -ldflags="-s -w" -o $(BINARY_PATH)/$(BINARY_NAME) $(CMD_DIR)

# Swagger
.PHONY: swag
swag:
	swag init --output $(SWAGGER_OUTPUT) --dir $(SWAGGER_DIRS) --generalInfo $(SWAGGER_GENERAL)

# Swagfmt
.PHONY: swagfmt
swagfmt:
	swag fmt --dir $(SWAGGER_DIRS)