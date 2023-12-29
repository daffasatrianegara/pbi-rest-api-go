package controllers

import (
    "net/http"
    "task-5-pbi-btpns-daffa_satria/app"
    "task-5-pbi-btpns-daffa_satria/config"
    "task-5-pbi-btpns-daffa_satria/helpers"
    "task-5-pbi-btpns-daffa_satria/models"
    "github.com/asaskevich/govalidator"
    "github.com/gin-gonic/gin"
)


func Register(context *gin.Context) {
    var userFormRegister app.Register
    if err := context.ShouldBindJSON(&userFormRegister); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        context.Abort()
        return
    }

    if _, err := govalidator.ValidateStruct(userFormRegister); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	
    var user models.User
    if err := config.GetDB().Where("email = ?", userFormRegister.Email).First(&user).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
        context.Abort()
        return
    }
	
	if len(userFormRegister.Password) < 6 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "password harus lebih dari 6 karakter"})
		context.Abort()
		return
	}

    user = models.User{
        Username: userFormRegister.Username,
        Email:    userFormRegister.Email,
        Password: userFormRegister.Password,
    }

    hashedPassword, err := helpers.HashPass(user.Password)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        context.Abort()
        return
    }

    user.Password = hashedPassword
    err = config.GetDB().Create(&user).Error
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        context.Abort()
        return
    }

    context.JSON(http.StatusCreated, gin.H{"message": "Register Berhasil"})
}



func Login(context *gin.Context) {
	var userFormLogin app.Login
	if err := context.ShouldBindJSON(&userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(userFormLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	if err := config.GetDB().Where("email = ?", userFormLogin.Email).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak ditemukan"})
		return
	}

	isMatch, err := helpers.CheckPass(userFormLogin.Password, user.Password)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isMatch {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Password salah"})
		return
	}

	token, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "token": token})
}