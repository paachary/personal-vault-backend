package apis

import (
	"go-mongo-project/config"
	"go-mongo-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPersonalDetails(context *gin.Context) {

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	var err error

	user, err = user.GetPersonalDetails()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"result":  user,
	})

}
