FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o academy .

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/web .

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/.env ./.env

EXPOSE 8080

ENTRYPOINT ["/app/web"]
