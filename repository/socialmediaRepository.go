package repository

import (
	"go-mygram/database"
	"go-mygram/models"
)

func FindAllSocialMedia(uid uint) ([]models.SocialMedia, error) {
	db := database.GetDB()
	SocialMedia := []models.SocialMedia{}

	err := db.Debug().Preload("User").Find(&SocialMedia, "user_id = ?", uid).Error
	return SocialMedia, err
}

func FindSocialMediaUser(id uint) (*models.SocialMedia, error) {
	db := database.GetDB()
	SocialMedia := models.SocialMedia{}

	err := db.Debug().Select("user_id").First(&SocialMedia, id).Error
	return &SocialMedia, err
}

func FindSocialMediaById(id uint) (*models.SocialMedia, error) {
	db := database.GetDB()
	SocialMedia := models.SocialMedia{}

	err := db.Debug().Preload("User").First(&SocialMedia, id).Error
	return &SocialMedia, err
}

func CreateSocialMedia(SocialMedia *models.SocialMedia) error {
	db := database.GetDB()

	err := db.Debug().Create(&SocialMedia).Error
	return err
}

func UpdateSocialMedia(SocialMediainput *models.SocialMedia) (*models.SocialMedia, error) {
	db := database.GetDB()

	SocialMedia, _ := FindSocialMediaById(SocialMediainput.ID)
	err := db.Model(&SocialMedia).Updates(&SocialMediainput).Error

	return SocialMedia, err
}

func DeleteSocialMedia(id uint) error {
	db := database.GetDB()
	SocialMedia := models.SocialMedia{}

	err := db.Debug().Where("id = ?", id).Delete(&SocialMedia).Error

	return err
}