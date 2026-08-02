package dto;

type Response struct{
	Id uint `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	CreatedAt string `json:"created_at"`
}