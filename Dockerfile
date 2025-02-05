FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /rate-limiter cmd/main.go

EXPOSE 8080

CMD ["/rate-limiter"]