package supabase

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	storage_go "github.com/supabase-community/storage-go"
)

type SupabaseStorageItf interface {
	Upload(bucket string, file *multipart.FileHeader) (string, error)
}

type SupabaseStorage struct {
	client *storage_go.Client
}

func NewSupabaseStorage(client *storage_go.Client) SupabaseStorageItf {
	return &SupabaseStorage{client: client}
}

func (s SupabaseStorage) Upload(bucket string, file *multipart.FileHeader) (string, error) {
	fileData, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileData.Close()
	fileReader := io.ReadCloser(fileData)
	defer fileReader.Close()

	fileName := filepath.Base(file.Filename)
	relativePath := fileName

	result, err := s.client.UploadFile(bucket, relativePath, fileReader,
		storage_go.FileOptions{ContentType: func() *string { s := "image/jpeg"; return &s }(),
			Upsert: func() *bool { b := true; return &b }()})
	if err != nil {
		return "", err
	}

	url := s.client.GetPublicUrl("", result.Key)

	cleanedURL := strings.Replace(url.SignedURL, "public//", "public/", -1)

	return cleanedURL, nil
}
