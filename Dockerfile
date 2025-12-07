FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

EXPOSE 8888

CMD ["./main"]