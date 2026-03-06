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

/** PRAN Related APIs **/

func AddPRANDetails(context *gin.Context) {

	pranDetails := models.Pran_Details{}

	if err := context.ShouldBindJSON(&pranDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	pranDetailsBson := bson.M{
		config.PRAN_NUMBER:   pranDetails.Pran_Number,
		config.PRAN_PASSWORD: pranDetails.Password,
	}

	encryptedDetails, err := utils.EncryptStructFields(pranDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	id := utils.GenerateUniqueKey(pranDetails.Pran_Number)

	encryptedDetails[config.PRAN_ID] = id

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.PRAN_DETAILS, config.PRAN_ID): bson.M{"$ne": id},
	}

	err = user.AddRecord(keyFilter, config.PRAN_DETAILS, encryptedDetails, filter, true)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "PRAN Details added to the document successfully.",
	})

}

func UpdatePRANDetails(context *gin.Context) {

	pranDetails := models.Pran_Details{}

	if err := context.ShouldBindJSON(&pranDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{
		fmt.Sprintf("elem.%s", config.PRAN_ID): pranDetails.Id}

	data := map[string]any{
		config.PRAN_PASSWORD: pranDetails.Password,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.PRAN_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "PRAN Details updated in the document successfully.",
	})

}

func DeletePRANDetails(context *gin.Context) {

	pranDetails := models.Pran_Details{}

	if err := context.ShouldBindJSON(&pranDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.PRAN_ID: pranDetails.Id}

	err := user.DeleteRecord(keyFilter, config.PRAN_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "PRAN Details deleted from the document successfully.",
	})

}

/*** End PRAN Related APIs **/
