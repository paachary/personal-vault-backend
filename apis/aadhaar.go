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

/** Aadhaar Related APIs **/

func AddAadhaarDetails(context *gin.Context) {

	aadhaarDetails := models.Aadhaar_Details{}

	if err := context.ShouldBindJSON(&aadhaarDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	aadhaarDetailsBson := bson.M{
		config.AADHAAR_NUMBER:     aadhaarDetails.Aadhaar_Number,
		config.AADHAAR_ISSUE_DATE: aadhaarDetails.Issue_Date,
	}

	encryptedDetails, err := utils.EncryptStructFields(aadhaarDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	encryptedDetails[config.AADHAAR_ID] = utils.GenerateUUID()

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.AADHAAR_DETAILS, config.AADHAAR_ID): bson.M{"$ne": aadhaarDetails.Id},
	}

	err = user.AddRecord(keyFilter, config.AADHAAR_DETAILS, encryptedDetails, filter, true)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "Aadhaar Details added to the document successfully.",
	})

}

func UpdateAadhaarDetails(context *gin.Context) {

	aadhaarDetails := models.Aadhaar_Details{}

	if err := context.ShouldBindJSON(&aadhaarDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{fmt.Sprintf("elem.%s", config.AADHAAR_ID): aadhaarDetails.Id}

	data := map[string]any{
		config.AADHAAR_ISSUE_DATE: aadhaarDetails.Issue_Date,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.AADHAAR_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Aadhaar Details updated in the document successfully.",
	})

}

func DeleteAadhaarDetails(context *gin.Context) {

	aadhaarDetails := models.Aadhaar_Details{}

	if err := context.ShouldBindJSON(&aadhaarDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.AADHAAR_ID: aadhaarDetails.Id}

	err := user.DeleteRecord(keyFilter, config.AADHAAR_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Aadhaar Details deleted from the document successfully.",
	})

}

/*** End Aadhaar Related APIs **/
