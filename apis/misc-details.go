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

/** Misc Related APIs **/

func AddMiscDetails(context *gin.Context) {

	miscDetails := models.Misc_Details{}

	if err := context.ShouldBindJSON(&miscDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	miscDetailsBson := bson.M{
		config.MISC_TYPE_CODE:   miscDetails.Type_Code,
		config.MISC_DESCRIPTION: miscDetails.Description,
		config.MISC_KEY_1:       miscDetails.Key_1,
		config.MISC_VAL_1:       miscDetails.Val_1,
		config.MISC_KEY_2:       miscDetails.Key_2,
		config.MISC_VAL_2:       miscDetails.Val_2,
	}

	id := utils.GenerateUniqueKey(miscDetails.Type_Code)

	filter := bson.M{
		config.USER_NAME: user.User_Name,
		fmt.Sprintf("%s.%s", config.MISC_DETAILS, config.MISC_ID): bson.M{"$ne": id},
	}

	encryptedDetails, err := utils.EncryptStructFields(miscDetailsBson)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	encryptedDetails[config.MISC_ID] = id

	err = user.AddRecord(keyFilter, config.MISC_DETAILS, encryptedDetails, filter, false)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("%s Details added to the document successfully.", miscDetails.Type_Code),
	})

}

func UpdateMiscDetails(context *gin.Context) {

	miscDetails := models.Misc_Details{}

	if err := context.ShouldBindJSON(&miscDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{fmt.Sprintf("elem.%s", config.MISC_ID): miscDetails.Id}

	data := map[string]any{
		config.MISC_DESCRIPTION: miscDetails.Description,
		config.MISC_VAL_1:       miscDetails.Val_1,
		config.MISC_VAL_2:       miscDetails.Val_2,
	}

	encryptedData, err := utils.EncryptStructFields(data)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt details"})
		return
	}

	err = user.UpdateArrayRecord(keyFilter, config.MISC_DETAILS, filter, encryptedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s Details updated in the document successfully.", miscDetails.Type_Code),
	})

}

func DeleteMiscDetails(context *gin.Context) {

	miscDetails := models.Misc_Details{}

	if err := context.ShouldBindJSON(&miscDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *models.User = &models.User{}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	filter := bson.M{config.MISC_ID: miscDetails.Id}

	err := user.DeleteRecord(keyFilter, config.MISC_DETAILS, filter)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s Details deleted from the document successfully.", miscDetails.Type_Code),
	})

}

/*** End Misc Related APIs **/
