package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/media/dto"
	"ticketBooking/internal/utils/query"
)

type Service struct {
	repo Repository
	cld  *CloudinaryClient
	
}

func NewMediaService(repo Repository, cld *CloudinaryClient) *Service {
	return &Service{repo: repo, cld: cld}
}

// UploadFile handles the full flow:
// 1. upload to Cloudinary
// 2. save metadata in media table
// 3. return response that contains both DB id and Cloudinary public_id
func (s *Service) UploadFile(
	ctx context.Context,
	file multipart.File,
	header *multipart.FileHeader,
	req *dto.UploadRequest,
) (*dto.Response, error) {

	// 1. Cloudinary upload
	uploadRes, err := s.cld.Upload(ctx, file, header, req.Folder, req.ModelName);

	fmt.Println("uploadRes ",uploadRes);
	
	if err != nil {
		return nil, err
	}

	// 2. Persist in DB
	media := &Media{
		ImageID:      uploadRes.PublicID,
		PublicURL:    uploadRes.URL,
		SecureURL:    uploadRes.SecureURL,
		ResourceType: uploadRes.ResourceType,
		Format:       uploadRes.Format,
		Type:         MapType(uploadRes.ResourceType, uploadRes.Format),
		ModelName:    req.ModelName,
		ModelID:      req.ModelID,
		Folder:       req.Folder,
		Size:         uploadRes.Bytes,
		Width:        uploadRes.Width,
		Height:       uploadRes.Height,
		OriginalName: header.Filename,
	}

	if media.Folder == "" {
		media.Folder = req.ModelName
	}

	if err := s.repo.Create(media); err != nil {
		// best-effort cleanup on Cloudinary if DB save fails
		_ = s.cld.Destroy(ctx, uploadRes.PublicID, uploadRes.ResourceType)
		return nil, fmt.Errorf("failed to save media record: %w", err)
	}

	return media.ToResponse(), nil
}






func (s *Service) GetByID(id uint) (*dto.Response, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return m.ToResponse(), nil
}

func (s *Service) List(params query.Params, filter dto.ListFilter) (*httpresponse.PaginatedData, error) {
	list, total, err := s.repo.List(params, filter)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.Response, 0, len(list))
	for _, m := range list {
		docs = append(docs, m.ToResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}

// DeleteByID:
// 1. find by DB auto id
// 2. destroy on Cloudinary using image_id (public_id)
// 3. soft-delete in DB
func (s *Service) DeleteByID(ctx context.Context, id uint) error {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("media not found: %w", err)
	}

	// Cloudinary destroy
	if err := s.cld.Destroy(ctx, m.ImageID, m.ResourceType); err != nil {
		// still hard-delete DB even if Cloudinary fails (log in real system)
		_ = s.repo.HardDelete(id)
		return fmt.Errorf("cloudinary destroy failed (db hard-deleted): %w", err)
	}

	return s.repo.HardDelete(id) 
}

// GetModelKeyed returns the shape the frontend asked for
// e.g. { "event_url": "...", "event_id": 12, "image_id": "..." }
func (s *Service) GetModelKeyed(id uint) (map[string]interface{}, error) {
	m, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return m.ToModelKeyedResponse(), nil
}
