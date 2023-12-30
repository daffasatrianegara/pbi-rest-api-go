package controllers

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"task-5-pbi-btpns-daffa_satria/app"
	"task-5-pbi-btpns-daffa_satria/config"
	"task-5-pbi-btpns-daffa_satria/helpers"
	"task-5-pbi-btpns-daffa_satria/models"
)

func GetAllPhotos(context *gin.Context) {
	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if ErrorHandling(context, err) {
		return
	}

	var photos []app.GetPhoto
	if err := config.GetDB().Table("photos").
		Select("photos.id, photos.title, photos.caption, photos.photo_url").
		Where("photos.user_id = ?", parsedToken.ID).
		Order("photos.id").
		Scan(&photos).
		Error; ErrorHandling(context, err) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": photos})
}

func GetPhotoById(context *gin.Context) {
	id := context.Param("id")
	tokenString := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(tokenString)
	if ErrorHandling(context, err) {
		return
	}

	var photo app.GetPhoto
	if err := config.GetDB().Table("photos").
		Select("photos.id, photos.title, photos.caption, photos.photo_url").
		Where("photos.id = ? AND photos.user_id = ?", id, parsedToken.ID).
		Scan(&photo).
		Error; ErrorHandling(context, err) {
		return
	}

	if photo.ID == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "Foto tidak ada"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"data": photo})
}


func UploadPhoto(context *gin.Context) {
	var photoUpload app.UploadPhoto
	if err := context.ShouldBindJSON(&photoUpload); ErrorHandling(context, err) {
		return
	}

	if _, err := govalidator.ValidateStruct(photoUpload); ErrorHandling(context, err) {
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if ErrorHandling(context, err) {
		return
	}

	photo := models.Photo{
		Title:    photoUpload.Title,
		Caption:  photoUpload.Caption,
		PhotoUrl: photoUpload.PhotoUrl,
		UserID:   parsedToken.ID,
	}

	if err := config.GetDB().Create(&photo).Error; ErrorHandling(context, err) {
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Foto berhasil terupload"})
}


func UpdatePhoto(context *gin.Context) {
	photoID, err := strconv.Atoi(context.Param("id"))
	if ErrorHandling(context, err) {
		return
	}

	var photoUpdate app.UpdatePhoto
	if err := context.ShouldBindJSON(&photoUpdate); ErrorHandling(context, err) {
		return
	}

	if _, err := govalidator.ValidateStruct(photoUpdate); ErrorHandling(context, err) {
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if ErrorHandling(context, err) {
		return
	}

	var photo models.Photo
	if err := config.GetDB().Where("id = ? AND user_id = ?", photoID, parsedToken.ID).First(&photo).Error; ErrorHandling(context, err) {
		return
	}

	photo.Title = photoUpdate.Title
	photo.Caption = photoUpdate.Caption
	photo.PhotoUrl = photoUpdate.PhotoUrl
	if err := config.GetDB().Save(&photo).Error; ErrorHandling(context, err) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Foto berhasil terupdate"})
}



func DeletePhoto(context *gin.Context) {
	var photo models.Photo
	photoID, err := strconv.Atoi(context.Param("id"))
	if ErrorHandling(context, err) {
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if ErrorHandling(context, err) {
		return
	}

	if err := config.GetDB().First(&photo, photoID).Error; ErrorHandling(context, err) {
		return
	}

	if photo.UserID != parsedToken.ID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	if err := config.GetDB().Delete(&photo).Error; ErrorHandling(context, err) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Foto berhasil dihapus"})
}


func ErrorHandling(context *gin.Context, err error) bool {
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return true
	}
	return false
}
