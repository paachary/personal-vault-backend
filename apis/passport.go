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

/** Passport Related APIs **/

func AddPassportDetails(context *gin.Context) {

	passportDetails := models.Passport_Details{}

	if err := context.ShouldBindJSON(&passportDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	passportDetailsBson := bson.M{
		config.PASSPORT_NUMBER:         passportDetails.Passport_Number,
		config.PASSPORT_ISSUE_DATE:     passportDetails.Issue_Date,
		config.PASSPORT_EXPIRY_DATE:    passportDetails.Expiry_Date,
		config.PASSPORT_ISSUER_COUNTRY: passportDetails.Issuer_Country,
		config.PASSPORT_USER_ID:        passportDetails.User_Id,
		config.PASSPORT_PASSWORD:       passportDetails.Password,
	}

	encryptedDetails, err := utils.EncryptStructFields(passportDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	id := utils.GenerateUniqueKey(passportDetails.Passport_Number)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.PASSPORT_DETAILS, config.PASSPORT_ID): bson.M{"$ne": id},
	}

	encryptedDetails[config.PASSPORT_ID] = id

	err = user.AddRecord(keyFilter, config.PASSPORT_DETAILS, encryptedDetails, filter, false)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "Passport Details added to the document successfully.",
	})

}

func UpdatePassportDetails(context *gin.Context) {

	passportDetails := models.Passport_Details{}

	if err := context.ShouldBindJSON(&passportDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{
		fmt.Sprintf("elem.%s", config.PASSPORT_ID): passportDetails.Id}

	data := map[string]any{
		config.PASSPORT_EXPIRY_DATE: passportDetails.Expiry_Date,
		config.PASSPORT_PASSWORD:    passportDetails.Password,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.PASSPORT_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Passport Details updated in the document successfully.",
	})

}

func DeletePassportDetails(context *gin.Context) {

	passportDetails := models.Passport_Details{}

	if err := context.ShouldBindJSON(&passportDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.PASSPORT_ID: passportDetails.Id}

	err := user.DeleteRecord(keyFilter, config.PASSPORT_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Passport Details deleted from the document successfully.",
	})

}

/*** End Passport Related APIs **/
