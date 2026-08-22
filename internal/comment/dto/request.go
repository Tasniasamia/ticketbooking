package dto

type CreateRequest struct {
	BlogID   uint   `json:"blog_id" validate:"required"`
	ParentID *uint  `json:"parent_id"` // null/omit = top-level comment; set = reply
	Content  string `json:"content" validate:"required,min=1"`
}

type UpdateRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}
