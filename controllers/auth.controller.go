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
	var RegisterUser app.Register
	if err := context.ShouldBindJSON(&RegisterUser); err != nil {
		handleError(context, http.StatusBadRequest, "Invalid JSON payload", err)
		return
	}

	if _, err := govalidator.ValidateStruct(RegisterUser); err != nil {
		handleError(context, http.StatusBadRequest, "Validation error", err)
		return
	}

	if err := checkExistingEmail(context, RegisterUser.Email); err != nil {
		return
	}

	if len(RegisterUser.Password) < 6 {
		handleError(context, http.StatusBadRequest, "Password must be at least 6 characters long", nil)
		return
	}

	user := models.User{
		Username: RegisterUser.Username,
		Email:    RegisterUser.Email,
		Password: RegisterUser.Password,
	}

	hashedPassword, err := helpers.HashPass(user.Password)
	if err != nil {
		handleError(context, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	user.Password = hashedPassword

	if err := config.GetDB().Create(&user).Error; err != nil {
		handleError(context, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Registration successful"})
}

func Login(context *gin.Context) {
	var LoginUser app.Login
	if err := context.ShouldBindJSON(&LoginUser); err != nil {
		handleError(context, http.StatusBadRequest, "Invalid JSON payload", err)
		return
	}

	if _, err := govalidator.ValidateStruct(LoginUser); err != nil {
		handleError(context, http.StatusBadRequest, "Validation error", err)
		return
	}

	user, err := getUserByEmail(context, LoginUser.Email)
	if err != nil {
		return
	}

	isMatch, err := helpers.CheckPass(LoginUser.Password, user.Password)
	if err != nil {
		handleError(context, http.StatusBadRequest, "Error checking password", err)
		return
	}
	if !isMatch {
		handleError(context, http.StatusBadRequest, "Incorrect password", nil)
		return
	}

	token, err := helpers.GenerateToken(user.ID, user.Email)
	if err != nil {
		handleError(context, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func handleError(context *gin.Context, statusCode int, message string, err error) {
	context.JSON(statusCode, gin.H{"error": message})
	if err != nil {
		context.Error(err)
	}
	context.Abort()
}

func checkExistingEmail(context *gin.Context, email string) error {
	var user models.User
	if err := config.GetDB().Where("email = ?", email).First(&user).Error; err == nil {
		handleError(context, http.StatusBadRequest, "Email already registered", nil)
		return err
	}
	return nil
}

func getUserByEmail(context *gin.Context, email string) (models.User, error) {
	var user models.User
	if err := config.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		handleError(context, http.StatusBadRequest, "Email not found", err)
		return user, err
	}
	return user, nil
}
