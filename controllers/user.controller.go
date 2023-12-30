package controllers

import (
	"net/http"
	"strconv"
	"task-5-pbi-btpns-daffa_satria/app"
	"task-5-pbi-btpns-daffa_satria/config"
	"task-5-pbi-btpns-daffa_satria/helpers"
	"task-5-pbi-btpns-daffa_satria/models"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)


func GetUserById(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if ErrorHandling(context, err) {
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if ErrorHandling(context, err) {
		return
	}

	if userID != int(parsedToken.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak sesuai"})
		context.Abort()
		return
	}

	var user models.User
	if err := config.GetDB().Where("id = ?", parsedToken.ID).First(&user).Error; ErrorHandling(context, err) {
		return
	}

	if err := config.GetDB().First(&user, userID).Error; ErrorHandling(context, err) {
		return
	}

	getUser := app.GetUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	context.JSON(http.StatusOK, gin.H{"data": getUser})
}


func UpdateUser(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if ErrorHandling(context, err) {
		return
	}

	var FormUpdate app.UpdateUser
	if err := context.ShouldBindJSON(&FormUpdate); ErrorHandling(context, err) {
		return
	}

	if _, err := govalidator.ValidateStruct(FormUpdate); ErrorHandling(context, err) {
		return
	}

	var user models.User
	if len(FormUpdate.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "password harus lebih dari 6 karakter"})
		context.Abort()
		return
	}

	if err := validateDuplicateEmail(context, FormUpdate.Email, userID); err != nil {
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)
	if ErrorHandling(context, err) {
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak sesuai"})
		context.Abort()
		return
	}

	if err := config.GetDB().Where("id = ?", claims.ID).First(&user).Error; ErrorHandling(context, err) {
		return
	}

	if err := config.GetDB().First(&user, userID).Error; ErrorHandling(context, err) {
		return
	}

	user.Username = FormUpdate.Username
	user.Email = FormUpdate.Email

	if FormUpdate.Password != "" {
		hashedPassword, err := helpers.HashPass(FormUpdate.Password)
		if ErrorHandling(context, err) {
			return
		}
		user.Password = hashedPassword
	}

	if err := config.GetDB().Save(&user).Error; ErrorHandling(context, err) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Data pengguna berhasil terupdate"})
}


func DeleteUser(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if ErrorHandling(context, err) {
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if ErrorHandling(context, err) {
		return
	}

	if userID != int(parsedToken.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token null"})
		context.Abort()
		return
	}

	var user models.User
	if err := config.GetDB().Where("id = ?", parsedToken.ID).First(&user).Error; ErrorHandling(context, err) {
		return
	}

	if err := config.GetDB().First(&user, userID).Error; ErrorHandling(context, err) {
		return
	}

	if err := config.GetDB().Delete(&user).Error; ErrorHandling(context, err) {
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Data pengguna berhasil terhapus"})
}

func validateDuplicateEmail(context *gin.Context, email string, userID int) error {
	var user models.User
	if err := config.GetDB().Where("email = ? AND id != ?", email, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		context.Abort()
		return err
	}
	return nil
}
