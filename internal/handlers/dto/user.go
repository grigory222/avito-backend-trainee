package dto

const (
	RoleEmployee  = "employee"
	RoleModerator = "moderator"
)

type RoleStruct struct {
	Role string `json:"role"`
}

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserCreateResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
