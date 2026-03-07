package entities

type UserEntity struct {
	ID       string
	Username string
	Email    string
	Password string
}

func NewUserEntity(id, username, email, password string) UserEntity {
	return UserEntity{
		ID:       id,
		Username: username,
		Email:    email,
		Password: password,
	}
}
