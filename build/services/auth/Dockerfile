FROM golang:1.23-alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o auth ./cmd/auth/main.go


FROM alpine

WORKDIR /build

COPY --from=builder /build/auth .

ENTRYPOINT ["./auth"]

EXPOSE 8000