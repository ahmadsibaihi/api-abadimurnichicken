FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base ca-certificates tzdata

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -mod=vendor -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]