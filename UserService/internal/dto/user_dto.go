package dto

type UserLoginRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}
type VerifyTokenRequest struct {
	Token string `json:"token" form:"token" validate:"required,max=255"`
}
type UserCreateRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email,max=255"`
	Name     string `json:"name" form:"name" validate:"required,min=4,max=255"`
	Password string `json:"password" form:"password" validate:"required,min=8,max=255"`
}

type UserupdatePasswordRequest struct {
	Password string `json:"password" form:"password" validate:"required,min=8,max=255"`
}

type UserupdateEmailRequest struct {
	Email string `json:"email" form:"email" validate:"required,email,max=255"`
}

type UserUpdateNameProfileRequest struct {
	Name string `json:"name" form:"name" validate:"required,min=4,max=255"`
}

type UserCreateResponse struct {
	Id int `json:"id"`
}

type UserResponse struct {
	Id    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
