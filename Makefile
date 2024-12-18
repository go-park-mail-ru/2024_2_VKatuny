# define executable
EXECUTABLE=uArt

# source dir
SRC_DIR=./cmd/app

# flags for compilation
BUILD_FLAGS=

# domain dirs in internal/pkg
DOMAINS=applicant cvs employer portfolio session vacancies file_loading
MICROSERVICE_DOMAINS = notifications



.PHONY: mock-gen api

all: install build

build:
	go build $(BUILD_FLAGS) -o $(EXECUTABLE) $(SRC_DIR)

install:
	go mod tidy

tests:
	go test ./... -coverprofile=coverage.out.tmp
	cat coverage.out.tmp | grep -v -E 'docs|mock|pb.go|_easyjson.go' > coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o index.html

clean:
	go clean

lint:
	revive -config reviveconfig.toml -formatter friendly ./...

api:
	swag init --generalInfo ./cmd/app/main.go --output api/ --pd
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
	@for domain in $(MICROSERVICE_DOMAINS); do \
		echo "Generating mocks for domain: $$domain"; \
		rm -rf microservice/$$domain/mock; \
		mockgen -source=microservices/$$domain/$$domain/$$domain.go \
		    -destination=microservices/$$domain/$$domain/mock/$$domain.go \
			-package=mock; \
	done
	@echo "Generating mocks for auth microservice"
	@rm -rf microservices/auth/mock
	@mockgen -source=microservices/auth/gen/auth_grpc.pb.go -destination=microservices/auth/mock/mock_grpc.go -package=mock
	@mockgen -source=microservices/auth/auth.go -destination=microservices/auth/mock/mock_delivery.go -package=mock
	@echo "Generating gRPC mocks for notifications microservice"
	@rm -rf microservices/notifications/mock
	@mockgen -source=microservices/notifications/generated/notifications_grpc.pb.go -destination=microservices/notifications/mock/mock_grpc.go -package=mock
	@echo "OK!"

redis-start:
	redis-server .\configs\redis.conf

auth-microservice:
	go run ./cmd/auth/main.go

run:
	go run $(SRC_DIR)/main.go

easyjson-gen:
	easyjson -all ./internal/pkg/dto
	easyjson -all ./microservices/notifications/notifications/dto