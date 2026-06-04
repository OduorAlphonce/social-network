# Build stage
FROM golang:1.22-alpine AS builder

# Install build-base (gcc, make, etc.) for CGO (sqlite3)
RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/pkg/db/migrations/sqlite ./pkg/db/migrations/sqlite

EXPOSE 8080

CMD ["./server"]
