package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"ticketBooking/internal/config"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryClient struct {
	cld *cloudinary.Cloudinary
	config config.Config
}

func NewCloudinaryClient(config config.Config) (*CloudinaryClient, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials missing (CLOUDINARY_CLOUD_NAME / API_KEY / API_SECRET)")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to init cloudinary: %w", err)
	}
	return &CloudinaryClient{cld: cld, config: config}, nil
}

// UploadResult is the minimal data we need after a successful upload
type UploadResult struct {
	PublicID     string
	URL          string
	SecureURL    string
	ResourceType string
	Format       string
	Bytes        int64
	Width        int
	Height       int
}

func (c *CloudinaryClient) Upload(
	ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
	folder string,
	modelName string,
) (*UploadResult, error) {

	if folder == "" {
		folder = modelName // default folder = model_name (event, user, ...)
	}

	// public_id = folder/timestamp_originalname (without extension)
	ext := filepath.Ext(header.Filename)
	base := strings.TrimSuffix(header.Filename, ext)
	base = strings.ReplaceAll(base, " ", "_")
	publicID := fmt.Sprintf("%d_%s",time.Now().UnixNano(), base)

	params := uploader.UploadParams{
		PublicID:     publicID,
		Folder:       folder,
		ResourceType: "auto",
		Overwrite:    boolPtr(false),
	}

res, err := c.cld.Upload.Upload(ctx, file, params)
if err != nil {
	return nil, fmt.Errorf("cloudinary upload failed: %w", err)
}

// ✅ empty result হলে error দাও
if res.PublicID == "" || res.SecureURL == "" {
	return nil, fmt.Errorf("cloudinary upload returned empty result: %+v", res.Error)
}

fmt.Println("OK public_id:", res.PublicID)
fmt.Println("OK secure_url:", res.SecureURL)

return &UploadResult{
	PublicID:     res.PublicID,
	URL:          res.URL,
	SecureURL:    res.SecureURL,
	ResourceType: res.ResourceType,
	Format:       res.Format,
	Bytes:        int64(res.Bytes),
	Width:        res.Width,
	Height:       res.Height,
}, nil
}

func (c *CloudinaryClient) Destroy(ctx context.Context, publicID, resourceType string) error {
	if resourceType == "" {
		resourceType = "image"
	}
	_, err := c.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})
	return err
}

// Map Cloudinary resource_type + format → our friendly "type"
func MapType(resourceType, format string) string {
	format = strings.ToLower(format)
	switch resourceType {
	case "image":
		return "image"
	case "video":
		return "video"
	case "raw":
		if format == "pdf" {
			return "pdf"
		}
		return "file"
	default:
		return "file"
	}
}

func boolPtr(b bool) *bool { return &b }
