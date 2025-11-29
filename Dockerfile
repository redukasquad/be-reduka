FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

COPY .env .env

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]