FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/eyeOne ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/eyeOne .

# Copy .env file if you want to provide config via file
# COPY .env .

EXPOSE 8080

CMD ["./eyeOne"]
