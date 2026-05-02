# Pakai versi Go yang sesuai (di gambar kamu pakai 1.21+)
FROM golang:1.21-alpine

# Install git karena library kamu (fiber, gorm) butuh ini buat download
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file kodingan
COPY . .

# --- BAGIAN PENTING ---
# Kita arahkan build ke folder cmd/api/main.go sesuai gambar kamu
RUN go build -o main ./cmd/api/main.go

EXPOSE 8080

CMD ["./main"]
