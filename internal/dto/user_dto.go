package dto

type CreateUserRequest struct {
	Name      string `form:"name"`
	Birthdate string `form:"birthdate"`
	// Image file is handled separately in the controller
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	ImageUrl string `json:"image_url"`
}
