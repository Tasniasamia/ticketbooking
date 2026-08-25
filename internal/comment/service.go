package comment

import (
	"errors"
	"ticketBooking/internal/comment/dto"
)

// BlogCounter — blog.comment_count atomic update (interface → no circular import)
type BlogCounter interface {
	IncrementCommentCount(blogID uint) error
	DecrementCommentCount(blogID uint) error
}

type Service struct {
	repo        Repository
	blogCounter BlogCounter // optional; nil হলে count skip
}

func NewCommentService(repo Repository, blogCounter BlogCounter) *Service {
	return &Service{repo: repo, blogCounter: blogCounter}
}

func (s *Service) CreateComment(userID uint, req *dto.CreateRequest) (*dto.Response, error) {
	if req.ParentID != nil {
		parent, err := s.repo.GetByID(*req.ParentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}
		if parent.BlogID != req.BlogID {
			return nil, errors.New("parent comment does not belong to this blog")
		}
	}

	c := &Comment{
		BlogID:   req.BlogID,
		UserID:   userID,
		ParentID: req.ParentID,
		Content:  req.Content,
	}

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}

	if s.blogCounter != nil {
		_ = s.blogCounter.IncrementCommentCount(req.BlogID)
	}

	return c.ToResponse(), nil
}

func (s *Service) GetCommentsTree(blogID uint) ([]*dto.Response, error) {
	flat, err := s.repo.GetByBlogID(blogID)
	if err != nil {
		return nil, err
	}
	return buildTree(flat), nil
}

func (s *Service) UpdateComment(commentID, userID uint, req *dto.UpdateRequest, isAdmin bool) (*dto.Response, error) {
	c, err := s.repo.GetByID(commentID)
	if err != nil {
		return nil, err
	}
	if !isAdmin && c.UserID != userID {
		return nil, errors.New("you can only edit your own comment")
	}
	if c.IsDeleted {
		return nil, errors.New("cannot edit a deleted comment")
	}

	c.Content = req.Content
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c.ToResponse(), nil
}

func (s *Service) DeleteComment(commentID, userID uint, isAdmin bool) error {
	c, err := s.repo.GetByID(commentID)
	if err != nil {
		return err
	}
	if !isAdmin && c.UserID != userID {
		return errors.New("you can only delete your own comment")
	}

	if err := s.repo.SoftDelete(commentID); err != nil {
		return err
	}

	if s.blogCounter != nil {
		_ = s.blogCounter.DecrementCommentCount(c.BlogID)
	}
	return nil
}

func buildTree(comments []*Comment) []*dto.Response {
	if len(comments) == 0 {
		return []*dto.Response{}
	}

	// [1] {Id, blogId,content}

	nodeMap := make(map[uint]*dto.Response)  
	var roots []*dto.Response

	for _, c := range comments {
		nodeMap[c.ID] = c.ToResponse()
	}

	for _, c := range comments {
		node := nodeMap[c.ID]
		if c.ParentID == nil {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[*c.ParentID]; ok {
				parent.Replies = append(parent.Replies, node)
			} else {
				roots = append(roots, node)
			}
		}
	}

	return roots
}
