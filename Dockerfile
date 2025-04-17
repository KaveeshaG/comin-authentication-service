FROM golang:1.23.4-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git

# Set GO111MODULE to on to ensure module mode is used
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Direct build command
RUN go build -o main ./cmd/auth

FROM alpine:3.18
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/main /app/auth-service
COPY --from=builder /app/.air.toml /app/ 2>/dev/null || true
COPY --from=builder /app/.env /app/ 2>/dev/null || true
ENV GO_ENV=production
EXPOSE 8080
CMD ["/app/auth-service"]
