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

/** PAN Related APIs **/

func AddPANDetails(context *gin.Context) {

	panDetails := models.Pan_Details{}

	if err := context.ShouldBindJSON(&panDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	panDetailsBson := bson.M{
		config.PAN_NUMBER:     panDetails.Pan_Number,
		config.PAN_ISSUE_DATE: panDetails.Issue_Date,
		config.PAN_PASSWORD:   panDetails.Password,
	}

	encryptedDetails, err := utils.EncryptStructFields(panDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	id := utils.GenerateUniqueKey(panDetails.Pan_Number)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.PAN_DETAILS, config.PAN_ID): bson.M{"$ne": id},
	}

	encryptedDetails[config.PAN_ID] = id

	err = user.AddRecord(keyFilter, config.PAN_DETAILS, encryptedDetails, filter, true)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "PAN Details added to the document successfully.",
	})

}

func UpdatePANDetails(context *gin.Context) {

	panDetails := models.Pan_Details{}

	if err := context.ShouldBindJSON(&panDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{
		fmt.Sprintf("elem.%s", config.PAN_ID): panDetails.Id}

	data := map[string]any{
		config.PAN_ISSUE_DATE: panDetails.Issue_Date,
		config.PAN_PASSWORD:   panDetails.Password,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.PAN_DETAILS, filter, encryptedData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "PAN Details updated in the document successfully.",
	})

}

func DeletePANDetails(context *gin.Context) {

	panDetails := models.Pan_Details{}

	if err := context.ShouldBindJSON(&panDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.PAN_ID: panDetails.Id}

	err := user.DeleteRecord(keyFilter, config.PAN_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "PAN Details deleted from the document successfully.",
	})

}

/*** End PAN Related APIs **/
