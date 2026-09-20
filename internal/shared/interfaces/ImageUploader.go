package interfaces

import "io"

type IImageUploader interface {
	Upload(file io.Reader, folderName string) (secureURL string, photoID string, err error)
}
