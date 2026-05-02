# 1. Gunakan image Golang versi terbaru yang ringan
FROM golang:1.21-alpine

# 2. Tentukan folder kerja di dalam server
WORKDIR /app

# 3. Copy file library (dependency)
COPY go.mod ./
# Jika kamu punya go.sum, aktifkan baris di bawah ini (hapus tanda #)
# COPY go.sum ./

# 4. Download library yang dibutuhkan
RUN go mod download

# 5. Copy seluruh sisa kodingan kamu
COPY . .

# 6. Build aplikasi kamu jadi file bernama 'main'
RUN go build -o main .

# 7. Beritahu Docker port mana yang dibuka (sesuaikan port Go kamu)
EXPOSE 8080

# 8. Jalankan aplikasinya
CMD ["./main"]
