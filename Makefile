# define executable
EXECUTABLE=uArt

# siurce dir
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