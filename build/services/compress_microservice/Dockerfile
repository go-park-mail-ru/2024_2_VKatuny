FROM golang:1.23-bullseye AS builder

RUN apt-get update -y && apt-get install -y pkg-config && apt-get install -y libvips-dev

WORKDIR /build

# ADD go.mod .
COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=1 go build -o pic-compress ./cmd/compress_microservice/main.go


FROM ubuntu

RUN apt-get update -y && apt-get install -y pkg-config && apt-get install -y libvips-dev

WORKDIR /build

COPY --from=builder /build/pic-compress .

ENTRYPOINT ["./pic-compress"]

EXPOSE 8000
