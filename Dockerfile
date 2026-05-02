FROM golang:1.21-alpine

# Set folder kerja
WORKDIR /app

# Copy seluruh kodingan (termasuk folder vendor tadi)
COPY . .

# Build aplikasi menggunakan folder vendor (-mod=vendor)
# Arahkan ke lokasi main.go kamu di cmd/api/main.go
RUN go build -mod=vendor -o main ./cmd/api/main.go

EXPOSE 8080

CMD ["./main"]
