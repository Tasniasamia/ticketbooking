package dto

import userDto "ticketBooking/internal/user/dto"

import ("ticketBooking/internal/event")

type UserInfo struct {
	ID             uint         `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Role           userDto.RoleType `json:"role"`
	ProfileImage   string       `json:"profile_image,omitempty"`
	ProfileImageId uint         `json:"profile_image_id,omitempty"`
	PhoneNumber    string       `json:"phone_number,omitempty"`
	Country        string       `json:"country,omitempty"`
	Status         string       `json:"status,omitempty"`
}

type Response struct {
	ID          uint          `json:"id"`
	UserID      uint          `json:"user_id"`
	EventID     uint          `json:"event_id"`
	Quantity    int           `json:"quantity"`
	TotalPrice  float64       `json:"total_price"`
	Status      BookingStatus `json:"status"`
	BookingCode string        `json:"booking_code"`
	CreatedAt   string        `json:"created_at"`
	UserInfo     *UserInfo     `json:"user_info,omitempty"`
	EventInfo    *event.Event   `json:"event_info,omitempty"`
}




