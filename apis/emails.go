package apis

import (
	"fmt"
	"go-mongo-project/config"
	"go-mongo-project/models"
	"go-mongo-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/** Emails Related APIs **/

func AddEmailDetails(context *gin.Context) {

	emailDetails := models.Email{}

	if err := context.ShouldBindJSON(&emailDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	emailDetailsBson := bson.M{
		config.EMAIL_ENTITY:   emailDetails.Entity,
		config.EMAIL_ID:       emailDetails.Email_Id,
		config.EMAIL_PASSWORD: emailDetails.Password,
	}

	id := utils.GenerateUniqueKey(emailDetails.Email_Id, emailDetails.Entity)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		config.EMAIL_DETAILS: bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					config.EMAIL_CONSTANT_ID: id,
				},
			},
		},
	}

	encryptedDetails, err := utils.EncryptStructFields(emailDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	encryptedDetails[config.EMAIL_CONSTANT_ID] = id

	err = user.AddRecord(keyFilter, config.EMAIL_DETAILS, encryptedDetails, filter, false)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "Email Details added to the document successfully.",
	})

}

func UpdateEmailDetails(context *gin.Context) {

	emailDetails := models.Email{}

	if err := context.ShouldBindJSON(&emailDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{fmt.Sprintf("elem.%s", config.EMAIL_CONSTANT_ID): emailDetails.Id}

	data := map[string]any{
		config.EMAIL_PASSWORD: emailDetails.Password,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.EMAIL_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Email Details updated in the document successfully.",
	})

}

func DeleteEmailDetails(context *gin.Context) {

	emailDetails := models.Email{}

	if err := context.ShouldBindJSON(&emailDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.EMAIL_CONSTANT_ID: emailDetails.Id}

	err := user.DeleteRecord(keyFilter, config.EMAIL_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Email Details deleted from the document successfully.",
	})

}

/*** End Emails Related APIs **/
