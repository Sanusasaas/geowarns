FROM golang:1.24.8-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /geowarns cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /geowarns /app/geowarns
COPY --from=builder /app/.env .

COPY migrations migrations

RUN apk add --no-cache postgresql-client

EXPOSE 8080

CMD ["./geowarns"]
