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

/** IRCTC Related APIs **/

func AddIRCTCDetails(context *gin.Context) {

	irctcDetails := models.IRCTC_Details{}

	if err := context.ShouldBindJSON(&irctcDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	irctcDetailsBson := bson.M{
		config.IRCTC_USER_NAME: irctcDetails.User_Name,
		config.IRCTC_EMAIL_ID:  irctcDetails.Email_Id,
		config.IRCTC_PASSWORD:  irctcDetails.Password,
	}

	id := utils.GenerateUniqueKey(irctcDetails.User_Name)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.IRCTC_DETAILS, config.IRCTC_ID): bson.M{"$ne": id},
	}

	encryptedDetails, err := utils.EncryptStructFields(irctcDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	encryptedDetails[config.IRCTC_ID] = id

	err = user.AddRecord(keyFilter, config.IRCTC_DETAILS, encryptedDetails, filter, false)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "IRCTC added to the document successfully.",
	})

}

func UpdateIRCTCDetails(context *gin.Context) {

	irctcDetails := models.IRCTC_Details{}

	if err := context.ShouldBindJSON(&irctcDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{fmt.Sprintf("elem.%s", config.IRCTC_ID): irctcDetails.Id}

	data := map[string]any{
		config.IRCTC_PASSWORD: irctcDetails.Password,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.IRCTC_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "IRCTC updated in the document successfully.",
	})

}

func DeleteIRCTCDetails(context *gin.Context) {

	irctcDetails := models.IRCTC_Details{}

	if err := context.ShouldBindJSON(&irctcDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.IRCTC_ID: irctcDetails.Id}

	err := user.DeleteRecord(keyFilter, config.IRCTC_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "IRCTC deleted from the document successfully.",
	})

}

/*** End IRCTC Related APIs **/
