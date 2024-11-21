# define executable
EXECUTABLE=uArt

# source dir
SRC_DIR=./cmd/app

# flags for compilation
BUILD_FLAGS=

# domain dirs in internal/pkg
DOMAINS=applicant cvs employer portfolio session vacancies

.PHONY: mock-gen

all: install build

build:
	go build $(BUILD_FLAGS) -o $(EXECUTABLE) $(SRC_DIR)

install:
	go mod tidy

tests:
	go test ./... -coverprofile=coverage.out.tmp
	cat coverage.out.tmp | grep -v 'mock' > coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o index.html

clean:
	go clean

lint:
	revive -config reviveconfig.toml -formatter friendly ./...

api:
	swag init --parseInternal --pd --dir cmd/myapp/,delivery/handler/ --output api/
	node ./api/server.js

mock-gen:
	@echo "Generating mocks..."
	@for domain in $(DOMAINS); do \
		echo "Generating mocks for domain: $$domain"; \
		rm -rf internal/pkg/$$domain/mock; \
		mockgen -source=internal/pkg/$$domain/$$domain.go \
		    -destination=internal/pkg/$$domain/mock/$$domain.go \
			-package=mock; \
	done
	@echo "OK!"

run:
	go run $(SRC_DIR)/main.go
