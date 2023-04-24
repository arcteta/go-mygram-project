package repository

import (
	"go-mygram/database"
	"go-mygram/models"
)

func FindUser(username, email string) *models.User {

	db := database.GetDB()
	User_exist := models.User{}
	_ = db.Debug().Where("username = ?", username).Or("email = ?", email).First(&User_exist).Error

	return &User_exist
}

func FindUserByUsername(username string) (*models.User, error) {
	db := database.GetDB()
	User := models.User{}
	err := db.Debug().Where("username = ?", username).Take(&User).Error

	return &User, err
}

func CreateUser(user *models.User) error {
	db := database.GetDB()
	err := db.Debug().Create(&user).Error
	return err
}
