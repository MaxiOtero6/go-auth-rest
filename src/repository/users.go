package repository

import (
	"github.com/MaxiOtero6/go-auth-rest/database"
	"github.com/MaxiOtero6/go-auth-rest/model"
)

type IUserRepository interface {
	CreateUser(userData *model.BaseUser) model.User
	GetUser(username string) model.User
	GetUserWithPassword(emailOrUsername string) model.BaseUser
	GetAllUsers() []model.User
}

type UserRepository struct{}

func (repository *UserRepository) CreateUser(userData *model.BaseUser) model.User {
	var user model.User

	database.DB.Exec(`
			INSERT INTO users (username, email, password) 
			VALUES (?, ?, ?);
		`,
		userData.Username, userData.Email, userData.Password,
	)

	database.DB.Raw("SELECT username, email FROM users WHERE username = ?", userData.Username).Scan(&user)

	return user
}

func (repository *UserRepository) GetUser(username string) model.User {
	var user model.User

	database.DB.Raw("SELECT username, email FROM users WHERE username = ?", username).Scan(&user)

	return user
}

func (repository *UserRepository) GetUserWithPassword(emailOrUsername string) model.BaseUser {
	var user model.BaseUser

	database.DB.Raw("SELECT username, email, password FROM users WHERE username = ? OR email = ?", emailOrUsername, emailOrUsername).Scan(&user)

	return user
}

func (repository *UserRepository) GetAllUsers() []model.User {
	var users []model.User

	database.DB.Raw("SELECT username, email FROM users").Scan(&users)

	return users
}
