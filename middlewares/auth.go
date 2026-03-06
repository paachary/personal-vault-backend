package middlewares

import (
	"go-mongo-project/config"
	"go-mongo-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	// Authentication logic goes here
	bearerToken := context.Request.Header.Get(("Authorization"))

	if (bearerToken == "") || (len(bearerToken) < 7) || (bearerToken[0:7] != "Bearer ") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is Unauthorized to perform the operation!"})
		return
	}

	claims, err := utils.VerifyToken(bearerToken[7:])

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User is Unauthorized to perform the operation!"})
		return
	}
	var userName string

	if f, ok := claims[config.USER_NAME]; ok {
		userName = f.(string)
	}
	context.Set(config.USER_NAME, userName)
	context.Next()
}
