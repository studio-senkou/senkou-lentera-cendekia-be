package storage

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

func GenerateEncryptedFilename(originalFilename, prefix string) string {
	ext := filepath.Ext(originalFilename)

	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	hash := md5.New()
	hash.Write([]byte(originalFilename))
	hash.Write([]byte(time.Now().String()))
	hash.Write(randomBytes)

	encryptedFilename := hex.EncodeToString(hash.Sum(nil))

	timestamp := time.Now().Format("20060102150405")

	return fmt.Sprintf("%s_%s_%s%s", strings.ToUpper(prefix), strings.ToUpper(timestamp), encryptedFilename[:16], ext)
}

func IsValidImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}

	return false
}

func GetContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func UploadFileToStorage(file *multipart.FileHeader, storagePath, prefix string, folder *string) (string, error) {
	if folder == nil {
		folder = new(string)
		*folder = ""
	}

	originalFilename := file.Filename
	if !IsValidImageExtension(originalFilename) {
		return "", fmt.Errorf("invalid file type: %s", originalFilename)
	}

	encryptedFilename := GenerateEncryptedFilename(originalFilename, prefix)
	fullPath := filepath.Join(storagePath, encryptedFilename)

	ctx := context.Background()
	uploader := NewUploadService()
	uploadedPath, err := uploader.UploadFile(ctx, file, fullPath, folder)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return uploadedPath, nil
}

func RemoveFileFromStorage(path string) error {
	ctx := context.Background()
	uploader := NewUploadService()
	err := uploader.RemoveFile(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}
