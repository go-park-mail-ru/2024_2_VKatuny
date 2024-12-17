FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk update && apk add --no-cache curl tar 

RUN mkdir pdfium-linux-x64 && curl -L https://github.com/bblanchon/pdfium-binaries/releases/download/chromium%2F6886/pdfium-linux-x64.tgz | tar xz -C pdfium-linux-x64

ENV EXTRACT_PATH=/opt/pdfium    

ADD pdfium.pc .

RUN sed -i "s|{path}|${EXTRACT_PATH}|g" pdfium.pc

ADD go.mod .

COPY . .

RUN go build -o app ./cmd/app/main.go


FROM ubuntu

WORKDIR /build

RUN apt-get update -y && apt-get install -y xvfb libfontconfig wkhtmltopdf

COPY --from=builder /build/app .

COPY --from=builder /build/pdfium.pc /usr/lib/pkgconfig/

COPY --from=builder /build/pdfium-linux-x64/ ${EXTRACT_PATH}

COPY --from=builder /build/server.crt /build/server.key .

ENV LD_LIBRARY_PATH={EXTRACT_PATH}/lib

ENTRYPOINT ["./app"]

EXPOSE 8000
