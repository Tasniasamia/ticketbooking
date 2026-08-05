package translation

import (
	"ticketBooking/internal/translation/dto"
	"ticketBooking/internal/utils/i18n"

	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	Key string `json:"key" gorm:"type:varchar(255);uniqueIndex;not null"`
	// "values" reserved keyword — column নাম lang_values
	Values i18n.LocalizedString `json:"values" gorm:"column:lang_values;type:jsonb;not null"`
}

func (t *Translation) ToResponse() *dto.Response {
	return &dto.Response{
		ID:        t.ID,
		Key:       t.Key,
		Values:    t.Values,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}