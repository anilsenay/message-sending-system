FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd/server

EXPOSE 8080

CMD ["./service"]