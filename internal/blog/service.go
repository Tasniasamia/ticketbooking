package blog

import (
	"errors"
	"ticketBooking/internal/blog/dto"
	"ticketBooking/internal/httpresponse"
	"ticketBooking/internal/user"
	"ticketBooking/internal/utils/i18n"
	"ticketBooking/internal/utils/query"
)

// CommentTreeLoader — comment module থেকে tree আনার জন্য (circular import এড়াতে)
type CommentTreeLoader interface {
	GetCommentsTree(blogID uint) (interface{}, error)
}

type Service struct {
	repo         Repository
	userRepo     user.Repository
	commentLoader CommentTreeLoader
}

func NewBlogService(repo Repository, userRepo user.Repository) *Service {
	return &Service{repo: repo, userRepo: userRepo}
}

// SetCommentLoader — register করার সময় inject করবে
func (s *Service) SetCommentLoader(loader CommentTreeLoader) {
	s.commentLoader = loader
}

func (s *Service) CreateBlog(req *dto.CreateRequest) (*dto.RawResponse, error) {
	status := req.Status
	if status == "" {
		status = dto.Pending
	}

	blog := &Blog{
		Title:            i18n.LocalizedString(req.Title),
		ShortDescription: i18n.LocalizedString(req.ShortDescription),
		LongDescription:  i18n.LocalizedString(req.LongDescription),
		ThumbnailImage:   req.ThumbnailImage,
		Images:           req.Images,
		AuthorID:         req.AuthorID,
		Status:           status,
		Slug:             req.Slug,
	}

	if req.AuthorID != 0 {
		if _, err := s.userRepo.GetUserActiveById(req.AuthorID); err != nil {
			return nil, errors.New("invalid author ID")
		}
	}

	if err := s.repo.Create(blog); err != nil {
		return nil, err
	}
	return blog.ToRawResponse(), nil
}

func (s *Service) GetAllBlogs(params query.Params, lang string) (*httpresponse.PaginatedData, error) {
	blogs, total, err := s.repo.GetAll(params, lang)
	if err != nil {
		return nil, err
	}

	docs := make([]*dto.RawResponse, 0, len(blogs))
	for _, b := range blogs {
		docs = append(docs, b.ToRawResponse())
	}

	meta := httpresponse.BuildPaginationMeta(total, params.Page, params.Limit)
	return &httpresponse.PaginatedData{
		Docs:           docs,
		PaginationMeta: meta,
	}, nil
}

func (s *Service) GetBlogByID(blogID uint, lang string, userID uint) (*dto.DetailResponse, error) {
	blog, err := s.repo.GetByID(blogID)
	if err != nil {
		return nil, err
	}

	raw := blog.ToRawResponse()

	// current user liked কিনা
	if userID > 0 {
		liked, _ := s.repo.IsLiked(blogID, userID)
		raw.IsLiked = liked
	}

	detail := &dto.DetailResponse{
		RawResponse: raw,
		Comments:    []interface{}{},
	}

	// comments tree load
	if s.commentLoader != nil {
		tree, err := s.commentLoader.GetCommentsTree(blogID)
		if err == nil {
			detail.Comments = tree
		}
	}

	return detail, nil
}

func (s *Service) GetBlogByIDAdmin(blogID uint) (*dto.RawResponse, error) {
	blog, err := s.repo.GetByIDAdmin(blogID)
	if err != nil {
		return nil, err
	}
	return blog.ToRawResponse(), nil
}

func (s *Service) UpdateBlog(blogID uint, req *dto.UpdateRequest) (*dto.RawResponse, error) {
	blog, err := s.repo.GetByIDAdmin(blogID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		blog.Title = i18n.LocalizedString(req.Title)
	}
	if req.ShortDescription != nil {
		blog.ShortDescription = i18n.LocalizedString(req.ShortDescription)
	}
	if req.LongDescription != nil {
		blog.LongDescription = i18n.LocalizedString(req.LongDescription)
	}
	if req.ThumbnailImage.URL != "" || req.ThumbnailImage.ID != 0 {
		blog.ThumbnailImage = req.ThumbnailImage
	}
	if req.Images != nil {
		blog.Images = req.Images
	}
	if req.AuthorID != 0 {
		blog.AuthorID = req.AuthorID
	}
	if req.Status != "" {
		blog.Status = req.Status
	}
	if req.Slug != "" {
		blog.Slug = req.Slug
	}

	if err := s.repo.Update(blog); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetByIDAdmin(blogID)
	if err != nil {
		return nil, err
	}
	return updated.ToRawResponse(), nil
}

func (s *Service) DeleteBlog(blogID uint) error {
	blog, err := s.repo.GetByIDAdmin(blogID)
	if err != nil {
		return err
	}
	return s.repo.Delete(blog)
}

func (s *Service) UpdateBlogStatus(req *dto.UpdateBlogStatusRequest) error {
	blog, err := s.repo.GetByIDAdmin(req.BlogID)
	if err != nil {
		return err
	}
	blog.Status = req.Status
	return s.repo.Update(blog)
}

// ToggleLike — concurrent safe (DB unique + atomic counter)
func (s *Service) ToggleLike(blogID, userID uint) (*dto.LikeResponse, error) {
	liked, count, err := s.repo.ToggleLike(blogID, userID)
	if err != nil {
		return nil, err
	}
	return &dto.LikeResponse{
		Liked:     liked,
		LikeCount: count,
	}, nil
}
