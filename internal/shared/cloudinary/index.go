// Package cloudinaryfileuploader provides utilities for uploading files to Cloudinary.
//
// It wraps the cloudinary-go client to simplify uploading files and returning
// the secure URL and public ID of the uploaded resource.
package cloudinary

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CloudinaryCredientials() (*cloudinary.Cloudinary, context.Context) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, _ := cloudinary.New()
	cld.Config.URL.Secure = true
	ctx := context.Background()
	return cld, ctx
}

type CloudinaryUploader struct {
	Cld *cloudinary.Cloudinary
	Ctx context.Context
}

func NewCloudinaryInstance() *CloudinaryUploader {
	cloudinaryCredientials, cloudinaryCtx := CloudinaryCredientials()
	return &CloudinaryUploader{
		Cld: cloudinaryCredientials,
		Ctx: cloudinaryCtx,
	}
}

func (c *CloudinaryUploader) Upload(file io.Reader, folderName string) (string, string, error) {
	res, err := c.Cld.Upload.Upload(c.Ctx, file, uploader.UploadParams{
		Folder: folderName,
	})

	if err != nil {
		return "", "", err
	}

	return res.SecureURL, res.PublicID, nil
}
