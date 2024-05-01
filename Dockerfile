# syntax=docker/dockerfile:1.3
FROM golang:1.22 AS builder
WORKDIR /
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ ./
RUN make

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /netrunner-alt-gen /

ENTRYPOINT ["/netrunner-alt-gen"]
