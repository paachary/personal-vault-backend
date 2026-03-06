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

/** Bank Related APIs **/

func AddBankDetails(context *gin.Context) {

	bankDetails := models.Bank_Details{}

	if err := context.ShouldBindJSON(&bankDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	bankDetailsBson := bson.M{
		config.BANK_NAME:                   bankDetails.Bank_Name,
		config.BANK_CUSTOMER_ID:            bankDetails.Customer_Id,
		config.BANK_USER_ID:                bankDetails.User_Id,
		config.BANK_LOGIN_PASSWORD:         bankDetails.Login_Password,
		config.BANK_TRANSACTION_PASSWORD:   bankDetails.Transaction_Password,
		config.BANK_MOBILE_PIN:             bankDetails.Mobile_Login_Pin,
		config.BANK_MOBILE_TRANSACTION_PIN: bankDetails.Mobile_Transaction_Pin,
	}

	bankId := utils.GenerateUniqueKey(bankDetails.Bank_Name, bankDetails.User_Id, bankDetails.Customer_Id)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		config.BANK_DETAILS: bson.M{
			"$not": bson.M{
				"$elemMatch": bson.M{
					config.BANK_ID: bankId,
				},
			},
		},
	}

	encryptedDetails, err := utils.EncryptStructFields(bankDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	encryptedDetails[config.BANK_ID] = bankId

	err = user.AddRecord(keyFilter, config.BANK_DETAILS, encryptedDetails, filter, false)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "Bank Details added to the document successfully.",
	})

}

func UpdateBankDetails(context *gin.Context) {

	bankDetails := models.Bank_Details{}

	if err := context.ShouldBindJSON(&bankDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{fmt.Sprintf("elem.%s", config.BANK_ID): bankDetails.Id}

	data := map[string]any{
		config.BANK_LOGIN_PASSWORD:         bankDetails.Login_Password,
		config.BANK_TRANSACTION_PASSWORD:   bankDetails.Transaction_Password,
		config.BANK_MOBILE_PIN:             bankDetails.Mobile_Login_Pin,
		config.BANK_MOBILE_TRANSACTION_PIN: bankDetails.Mobile_Transaction_Pin,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.BANK_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Bank Details updated in the document successfully.",
	})

}

func DeleteBankDetails(context *gin.Context) {

	bankDetails := models.Bank_Details{}

	if err := context.ShouldBindJSON(&bankDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.BANK_ID: bankDetails.Id}

	err := user.DeleteRecord(keyFilter, config.BANK_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": "Bank Details deleted from the document successfully.",
	})

}

/*** End Bank Related APIs **/
