FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY main .

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]