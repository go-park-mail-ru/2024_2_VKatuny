# define executable
EXECUTABLE=uArt.exe

# source dir
SRC_DIR=./cmd/app

# flags for compilation
BUILD_FLAGS=

.PHONY: api

all: build

build:
	go build $(BUILD_FLAGS) -o $(EXECUTABLE) $(SRC_DIR)

tests:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out
	go tool cover -html cover.out -o index.html

clean:
	go clean

lint:
	revive -config reviveconfig.toml -formatter friendly ./...

api:
	swag init --generalInfo ./cmd/app/main.go --output api/
	node ./api/server.js

auth-microservice:
	go run ./cmd/auth/main.go

run:
	go run $(SRC_DIR)/main.go
