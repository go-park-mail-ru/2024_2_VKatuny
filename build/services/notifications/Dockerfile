FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum .

COPY . .

RUN go build -o notifications ./cmd/notifications_microservice/main.go


FROM alpine

WORKDIR /build

COPY --from=builder /build/notifications .

ENTRYPOINT ["./notifications"]  

EXPOSE 8000
