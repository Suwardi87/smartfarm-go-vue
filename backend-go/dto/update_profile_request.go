package dto

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required,min=1"`
	Email string `json:"email" binding:"required,email"`
}
