package user

import "github.com/tekramkcots/sdk/models"

type User struct {
	UserID      string
	UserName    string
	AccessToken string
	Email       string
}

func (u User) ToModel() models.User {
	return models.User{
		UserID:      u.UserID,
		UserName:    u.UserName,
		AccessToken: u.AccessToken,
		Email:       u.Email,
	}
}

func FromModel(u models.User) User {
	return User{
		UserID:   u.UserID,
		UserName: u.UserName,
		Email:    u.Email,
	}
}
