# define executable
EXECUTABLE=uArt.exe

# source dir
SRC_DIR=./cmd/app

# flags for compilation
BUILD_FLAGS=

all: build

build:
	go build $(BUILD_FLAGS) -o $(EXECUTABLE) $(SRC_DIR)

tests:
	go test -cover ./...

clean:
	go clean

lint:
	revive -config reviveconfig.toml -formatter friendly ./...

api:
	swag init --parseInternal --pd --dir cmd/myapp/,delivery/handler/ --output api/
	node ./api/server.js

run:
	go run $(SRC_DIR)/main.go
