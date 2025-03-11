FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/cmd/auth/main /app/auth-service

COPY --from=builder /app/.air.toml /app/
COPY --from=builder /app/.env /app/

ENV GO_ENV=production

EXPOSE 8080

CMD ["/app/auth-service"]