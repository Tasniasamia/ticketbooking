package dto

type CreateRequest struct {
	Key    string            `json:"key" validate:"required,min=1,max=255"`
	Values map[string]string `json:"values"`
}

type UpdateRequest struct {
	Values map[string]string `json:"values" validate:"required,min=1"`
}

type BulkUpdateRequest struct {
	Items []BulkItem `json:"items" validate:"required,min=1,dive"`
}

type BulkItem struct {
	Key    string            `json:"key" validate:"required,min=1,max=255"`
	Values map[string]string `json:"values" validate:"required,min=1"`
}