package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// SaveUploadedPhotos menyimpan foto (maksimal 3 file) ke direktori upload dan mengembalikan array path
func SaveUploadedPhotos(files []*multipart.FileHeader, uploadDir string, maxFileSizeBytes int64) ([]string, error) {
	if len(files) == 0 {
		return []string{}, nil
	}

	if len(files) > 3 {
		return nil, errors.New("maksimal foto yang diperbolehkan adalah 3 file")
	}

	// Pastikan direktori tujuan tersedia
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori upload: %w", err)
	}

	var savedPaths []string

	for i, fileHeader := range files {
		// Validasi ukuran
		if fileHeader.Size > maxFileSizeBytes {
			return nil, fmt.Errorf("file ke-%d (%s) melebihi batas ukuran maksimal %d MB", i+1, fileHeader.Filename, maxFileSizeBytes/(1024*1024))
		}

		// Validasi ekstensi
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedExtensions[ext] {
			return nil, fmt.Errorf("format file '%s' tidak diizinkan. Hanya menerima jpg, jpeg, png, webp", fileHeader.Filename)
		}

		// Generate nama berkas unik
		uniqueName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
		destination := filepath.Join(uploadDir, uniqueName)

		// Buka berkas sumber
		src, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("gagal membaca file '%s': %w", fileHeader.Filename, err)
		}
		defer src.Close()

		// Buat berkas tujuan
		dst, err := os.Create(destination)
		if err != nil {
			return nil, fmt.Errorf("gagal menyimpan file ke '%s': %w", destination, err)
		}
		defer dst.Close()

		if _, err := dst.ReadFrom(src); err != nil {
			return nil, fmt.Errorf("gagal menulis data file: %w", err)
		}

		// Simpan path relatif untuk URL publik
		relPath := fmt.Sprintf("/uploads/tickets/%s", uniqueName)
		savedPaths = append(savedPaths, relPath)
	}

	return savedPaths, nil
}
