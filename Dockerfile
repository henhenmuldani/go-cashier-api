# =========================
# Build stage
# =========================
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Copy go.mod dulu supaya cache Docker kepakai
COPY go.mod ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags netgo -ldflags "-s -w" -o app

# =========================
# Runtime stage
# =========================
FROM alpine:latest

WORKDIR /app

# Copy binary dari build stage
COPY --from=builder /app/app .

# Expose port (sesuaikan dengan main.go kamu)
EXPOSE 8080

# Run app
CMD ["./app"]
