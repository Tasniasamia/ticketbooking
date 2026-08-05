package dto

import "ticketBooking/internal/utils/i18n"

type Response struct {
	ID        uint                 `json:"id"`
	Key       string               `json:"key"`
	Values    i18n.LocalizedString `json:"values"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

// Frontend-friendly: key → { langCode: value }
type GroupResponse map[string]map[string]string