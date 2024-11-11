package usecase

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/federicodosantos/socialize/pkg/supabase"
)

type FileUsecaseItf interface {
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error)
}

type FileUsecase struct {
	supabase supabase.SupabaseStorageItf
}

func NewFileUsecase(supabase supabase.SupabaseStorageItf) FileUsecaseItf {
	return &FileUsecase{
		supabase: supabase,
	}
}

func (uc *FileUsecase) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	return uc.supabase.Upload(os.Getenv("SUPABASE_BUCKET_USER"), fileHeader)
}
