# Build stage
FROM golang:1.23.2 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o music cmd/main.go

# Deploy stage
FROM alpine:3.21

RUN apk add --no-cache libc6-compat

WORKDIR /app

COPY --from=builder /app/music .

RUN chmod +x /app/music

CMD ["./music"]