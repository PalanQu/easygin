package models

type User struct {
	ID   int
	Name string
}

type CreateUserRequest struct {
	Name string
}

type CreateUserResponse struct {
	ID int
}

type GetAllUsersResponse struct {
	Users []*User
}
