FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mikrotik-exporter .

FROM scratch
COPY --from=builder /build/mikrotik-exporter /mikrotik-exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 9090
ENTRYPOINT ["/mikrotik-exporter"]
