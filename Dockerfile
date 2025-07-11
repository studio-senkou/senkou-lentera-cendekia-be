FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]