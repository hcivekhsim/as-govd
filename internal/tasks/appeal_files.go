package tasks

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

var allowedExtensions = map[string]struct{}{
	".doc":  {},
	".docx": {},
	".txt":  {},
	".rtf":  {},
	".odf":  {},
	".pdf":  {},
	".jpg":  {},
	".jpeg": {},
	".tiff": {},
	".png":  {},
}

var allowedMIMEs = map[string]struct{}{
	"application/msword": {},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {},
	"text/plain":      {},
	"application/rtf": {},
	"text/rtf":        {},
	"application/vnd.oasis.opendocument.text": {},
	"application/pdf":                         {},
	"image/jpeg":                              {},
	"image/tiff":                              {},
	"image/png":                               {},
}

type SavedAppealFile struct {
	OriginalName string
	StoredName   string
	StoragePath  string
	ContentType  string
	SizeBytes    int64
}

func saveAppealFiles(headers []*multipart.FileHeader) ([]SavedAppealFile, error) {
	if len(headers) > runtimeSettings.MaxFilesCount {
		return nil, ErrRequestTooLarge
	}

	if err := os.MkdirAll(runtimeSettings.UploadDir, 0o755); err != nil {
		return nil, err
	}

	saved := make([]SavedAppealFile, 0, len(headers))
	for _, header := range headers {
		item, err := saveAppealFile(header)
		if err != nil {
			cleanupSavedFiles(saved)
			return nil, err
		}
		saved = append(saved, item)
	}

	return saved, nil
}

func saveAppealFile(header *multipart.FileHeader) (SavedAppealFile, error) {
	if header.Size > runtimeSettings.MaxFileSizeBytes {
		return SavedAppealFile{}, ErrFileTooLarge
	}

	originalName := filepath.Base(header.Filename)
	ext := strings.ToLower(filepath.Ext(originalName))
	if _, ok := allowedExtensions[ext]; !ok {
		return SavedAppealFile{}, ErrUnsupportedFileType
	}

	file, err := header.Open()
	if err != nil {
		return SavedAppealFile{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, runtimeSettings.MaxFileSizeBytes+1))
	if err != nil {
		return SavedAppealFile{}, err
	}
	if int64(len(data)) > runtimeSettings.MaxFileSizeBytes {
		return SavedAppealFile{}, ErrFileTooLarge
	}

	mtype := mimetype.Detect(data)
	contentType := mtype.String()
	if i := strings.Index(contentType, ";"); i != -1 {
		contentType = contentType[:i]
	}
	if _, ok := allowedMIMEs[contentType]; !ok {
		return SavedAppealFile{}, ErrUnsupportedFileType
	}

	storedName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randomSuffix(), ext)
	storagePath := filepath.Join(runtimeSettings.UploadDir, storedName)
	if err := os.WriteFile(storagePath, data, 0o600); err != nil {
		return SavedAppealFile{}, err
	}

	return SavedAppealFile{
		OriginalName: originalName,
		StoredName:   storedName,
		StoragePath:  storagePath,
		ContentType:  contentType,
		SizeBytes:    int64(len(data)),
	}, nil
}

func cleanupSavedFiles(files []SavedAppealFile) {
	for _, file := range files {
		_ = os.Remove(file.StoragePath)
	}
}

func randomSuffix() string {
	token, err := generateToken(8)
	if err != nil {
		return "fallback"
	}
	return strings.NewReplacer("/", "", "+", "", "=", "").Replace(token)
}
