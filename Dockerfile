# Stage 1: Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies yang dibutuhkan untuk kompilasi
RUN apk add --no-cache git ca-certificates tzdata

# Copy dependency definition
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Compile binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./cmd/api

# Stage 2: Runtime Stage (Ringan & Cepat)
FROM alpine:3.19

WORKDIR /app

# Pasang sertifikat SSL & data zona waktu
RUN apk add --no-cache ca-certificates tzdata

# Set zona waktu default ke WIB (Asia/Jakarta)
ENV TZ=Asia/Jakarta

# Salin binary dan dokumen Swagger
COPY --from=builder /app/server /app/server
COPY --from=builder /app/docs /app/docs

# Pastikan folder penyimpanan uploads tersedia
RUN mkdir -p /app/storage/uploads/tickets

# Port aplikasi
EXPOSE 8086

# Jalankan server
CMD ["/app/server"]
