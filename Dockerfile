FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o -trimpath -o /rate-limiter cmd/main.go

FROM scratch
COPY --from=builder /rate-limiter /rate-limiter
EXPOSE 8080
CMD ["/rate-limiter"]
