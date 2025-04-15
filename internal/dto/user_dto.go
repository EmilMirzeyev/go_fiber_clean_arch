package dto

type CreateUserRequest struct {
	Name      string `form:"name" validate:"required"`
	Email     string `form:"email" validate:"required,email"`
	Password  string `form:"password" validate:"required,min=6"`
	Birthdate string `form:"birthdate" validate:"required"`
	RoleName  string `form:"role_name" validate:"required,oneof=admin moderator user"`
	// Image file is handled separately in the controller
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	ImageUrl string `json:"image_url"`
	Role     string `json:"role"`
}
