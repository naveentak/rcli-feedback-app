FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /feedback-server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /feedback-server .
COPY web/ ./web/

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080
USER nobody

CMD ["./feedback-server"]