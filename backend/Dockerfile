# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Enable CGO off for smaller static-ish binary (optional)
ENV CGO_ENABLED=0 GO111MODULE=on

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN go build -o /app/server ./

FROM alpine:3.20
WORKDIR /app
RUN adduser -D -u 10001 appuser
COPY --from=builder /app/server /app/server

# app listens on 8000
EXPOSE 8000
USER appuser

ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=clocked \
    SECRET=replace_me

CMD ["/app/server"]

