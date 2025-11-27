FROM golang:1.24.2-alpine

WORKDIR /app

# Copy dan download dependencies
COPY go.* ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Copy .env (pastikan .dockerignore tidak mengabaikannya)
COPY .env .env

# Build binary
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
