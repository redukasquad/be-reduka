FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install certificates for HTTPS
RUN apk add --no-cache ca-certificates git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# Final stage - minimal image
FROM alpine:latest

WORKDIR /app

# Install certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Environment variables are set in Render Dashboard
# Do not copy .env file for security

EXPOSE 8888

CMD ["./main"]