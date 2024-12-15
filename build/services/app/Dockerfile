FROM golang:1.23-alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o app ./cmd/app/main.go


FROM alpine

WORKDIR /build

COPY --from=builder /build/app .

ENTRYPOINT ["./app"]

EXPOSE 8000