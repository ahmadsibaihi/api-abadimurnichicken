FROM golang:1.21-alpine
WORKDIR /app
# Kita pakai folder vendor yang sudah kamu buat tadi
COPY . .
# Perintah build harus mengarah ke tempat main.go berada
RUN go build -mod=vendor -o main ./cmd/api/main.go
EXPOSE 8080
CMD ["./main"]
