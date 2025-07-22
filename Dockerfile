# syntax=docker/dockerfile:1
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o queue ./cmd/queue/main.go


FROM alpine:3.20
RUN apk add --no-cache iproute2 iputils

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

USER appuser
WORKDIR /home/appuser

COPY --from=builder /app/queue /home/appuser/queue

CMD ["/home/appuser/queue"]
