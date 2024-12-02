# Build stage
FROM golang:1.23.2 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o go-app cmd/main.go

# Deploy stage
FROM alpine:latest

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /app/go-app .

RUN chmod +x /app/go-app

CMD ["./go-app"]