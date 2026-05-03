FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://proxy.golang.org,direct
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]