FROM golang:1.23-alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o pic-compress ./cmd/compress_microservice/main.go


FROM alpine

WORKDIR /build

COPY --from=builder /build/pic-compress .

ENTRYPOINT ["./pic-compress"]

EXPOSE 8000