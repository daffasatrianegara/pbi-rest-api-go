package middlewares

import (
	"github.com/gin-gonic/gin"
	"task-5-pbi-btpns-daffa_satria/helpers"
)

func Authentication() gin.HandlerFunc {
    return func(context *gin.Context) {
        tokenString := context.GetHeader("Authorization")
        if tokenString == "" {
            context.JSON(401, gin.H{"error": "Authorization header is missing"})
            context.Abort()
            return
        }
        err := helpers.ValidateToken(tokenString)
        if err != nil {
            context.JSON(401, gin.H{"error": err.Error()})
            context.Abort()
            return
        }
        context.Next()
    }
}