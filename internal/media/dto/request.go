package dto

// Upload comes as multipart/form-data from frontend.
// Fields that frontend should send:
//
//	file        → the actual file (required)
//	model_name  → "event" | "user" | "product" | ... (required)
//	model_id    → optional related entity primary key
//	folder      → optional Cloudinary folder override
type UploadRequest struct {
	ModelName string `form:"model_name" validate:"required,oneof=event user product booking category banner payments"`
	ModelID   *uint  `form:"model_id"`
	Folder    string `form:"folder"`
}

// Delete by database auto-generated ID
type DeleteRequest struct {
	ID uint `param:"id" validate:"required"`
}

// List / filter
type ListFilter struct {
	ModelName string
	ModelID   *uint
	Type      string // image | video | pdf | file
}
