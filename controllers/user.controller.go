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
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	if userID != int(parsedToken.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak sesuai"})
		context.Abort()
		return
	}
	var user models.User
	if err := config.GetDB().Where("id = ?", parsedToken.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := config.GetDB().First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	var getUser app.GetUser
	getUser.ID = user.ID
	getUser.Username = user.Username
	getUser.Email = user.Email

	context.JSON(http.StatusOK, gin.H{"data": getUser})
}



func UpdateUser(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "data pengguna tidak ditemukan"})
		return
	}
	var FormUpdate app.UpdateUser
	if err := context.ShouldBindJSON(&FormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(FormUpdate); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if len(FormUpdate.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "password harus lebih dari 6 karakter"})
		context.Abort()
		return
	}

	if err := config.GetDB().Where("email = ? AND id != ?", FormUpdate.Email, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		context.Abort()
		return
	}

	if err := config.GetDB().Where("username = ? AND id != ?", FormUpdate.Username, userID).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar"})
		context.Abort()
		return
	}

	tokenString := context.GetHeader("Authorization")
	claims, err := helpers.ParseToken(tokenString)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(claims.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak sesuai"})
		context.Abort()
		return
	}

	if err := config.GetDB().Where("id = ?", claims.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := config.GetDB().First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	user.Username = FormUpdate.Username
	user.Email = FormUpdate.Email
    if FormUpdate.Password != "" {
        hashedPassword, err := helpers.HashPass(FormUpdate.Password)
        if err != nil {
            context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            context.Abort()
            return
        }
        user.Password = hashedPassword
    }

	if err := config.GetDB().Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Data pengguna berhasil terupdate"})
}



func DeleteUser(context *gin.Context) {
	var user models.User
	userID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	token := context.GetHeader("Authorization")
	parsedToken, err := helpers.ParseToken(token)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if userID != int(parsedToken.ID) {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Token null"})
		context.Abort()
		return
	}

	if err := config.GetDB().Where("id = ?", parsedToken.ID).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := config.GetDB().First(&user, userID).Error; err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := config.GetDB().Delete(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Data pengguna berhasil terhapus"})
}
