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

/** Mutual Fund Related APIs **/

func AddMFDetails(context *gin.Context) {

	mfDetails := models.Mutual_Fund{}

	if err := context.ShouldBindJSON(&mfDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	mfDetailsBson := bson.M{
		config.MUTUAL_FUND_FOLIO_NUMBER:         mfDetails.Folio_Number,
		config.MUTUAL_FUND_NAME:                 mfDetails.Fund_Name,
		config.MUTUAL_FUND_USER_ID:              mfDetails.User_Id,
		config.MUTUAL_FUND_EMAIL_ID:             mfDetails.Email_Id,
		config.MUTUAL_FUND_MOBILE_NUMBER:        mfDetails.Mobile_Number,
		config.MUTUAL_FUND_LOGIN_PASSWORD:       mfDetails.Login_Password,
		config.MUTUAL_FUND_TRANSACTION_PASSWORD: mfDetails.Transaction_Password,
		config.MUTUAL_FUND_MPIN:                 mfDetails.Mpin,
	}

	id := utils.GenerateUniqueKey(mfDetails.Folio_Number)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.MUTUAL_FUND_DETAILS, config.MUTUAL_FUND_ID): bson.M{"$ne": id},
	}

	encryptedDetails, err := utils.EncryptStructFields(mfDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	encryptedDetails[config.MUTUAL_FUND_ID] = id

	err = user.AddRecord(keyFilter, config.MUTUAL_FUND_DETAILS, encryptedDetails, filter, false)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "Mutual Fund Details added to the document successfully.",
	})

}

func UpdateMFDetails(context *gin.Context) {

	mfDetails := models.Mutual_Fund{}

	if err := context.ShouldBindJSON(&mfDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{
		fmt.Sprintf("elem.%s", config.MUTUAL_FUND_ID): mfDetails.Id}

	data := map[string]any{
		config.MUTUAL_FUND_USER_ID:              mfDetails.User_Id,
		config.MUTUAL_FUND_LOGIN_PASSWORD:       mfDetails.Login_Password,
		config.MUTUAL_FUND_TRANSACTION_PASSWORD: mfDetails.Transaction_Password,
		config.MUTUAL_FUND_EMAIL_ID:             mfDetails.Email_Id,
		config.MUTUAL_FUND_MOBILE_NUMBER:        mfDetails.Mobile_Number,
		config.MUTUAL_FUND_MPIN:                 mfDetails.Mpin,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.MUTUAL_FUND_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Mutual Fund Details updated in the document successfully.",
	})
}

func DeleteMFDetails(context *gin.Context) {

	mfDetails := models.Mutual_Fund{}

	if err := context.ShouldBindJSON(&mfDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.MUTUAL_FUND_ID: mfDetails.Id}

	err := user.DeleteRecord(keyFilter, config.MUTUAL_FUND_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Mutual Fund Details deleted from the document successfully.",
	})

}

/*** End Mutual Fund Related APIs **/
