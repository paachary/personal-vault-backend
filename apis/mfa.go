package apis

import (
	"go-mongo-project/db"
	"go-mongo-project/models"
	"go-mongo-project/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateMFA(context *gin.Context) {
	input := models.MFA_Code{}

	// Bind JSON to user struct
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate MFA code with 5-minute expiration
	mfaCode := utils.GenerateMFA()

	mfaCodeDocument := bson.M{
		"user_name": input.User_Name,
		"email_id":  input.Email_Id,
		"code":      mfaCode.Code,
		"createdAt": time.Now(),
		"expiresAt": mfaCode.ExpiresAt,
		"verified":  false,
	}

	keyFilter := bson.M{
		"user_name": input.User_Name,
	}

	// Store mfaCode in database with user identifier
	err := db.UpsertMfaRecord(keyFilter, nil, mfaCodeDocument)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store MFA code"})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message":  "User logged in successfully",
		"mfa_code": mfaCode.Code,
	})
}

func VerifyMFA(context *gin.Context) {
	input := models.MFA_Code{}

	// Bind JSON to user struct
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	keyFilter := bson.M{
		"user_name": input.User_Name,
	}

	condition := bson.M{
		"user_name": input.User_Name,
		"email_id":  input.Email_Id,
		"code":      input.Code,
		"verified":  false}

	mfaRecord, err := db.GetMfaRecord(keyFilter, condition)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve MFA record"})
		return
	}

	// MongoDB _id is primitive.ObjectID, not string
	var id string
	switch v := (*mfaRecord)["_id"].(type) {
	case bson.ObjectID:
		id = v.Hex()
	case string:
		id = v
	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid MFA record ID format"})
		return
	}

	// Extract values from bson.M
	code, ok := (*mfaRecord)["code"].(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid MFA code format"})
		return
	}

	// Handle different time formats from MongoDB
	var expiresAt time.Time
	switch v := (*mfaRecord)["expiresAt"].(type) {
	case time.Time:
		expiresAt = v
	case bson.DateTime:
		expiresAt = v.Time()
	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid expiration time format"})
		return
	}

	storedMFA := utils.MFACode{
		Code:      code,
		ExpiresAt: expiresAt,
	}

	// Convert string ID back to ObjectID for MongoDB queries
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid ObjectID format"})
		return
	}

	condition = bson.M{
		"_id": objectID,
	}

	// Check if code is expired
	if storedMFA.IsExpired() {
		if err := db.DeleteMfaRecord(keyFilter, condition); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expired MFA record"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Verification Code has Expired"})
		return
	}

	// Verify the code matches
	if code != storedMFA.Code {
		if err := db.DeleteMfaRecord(keyFilter, condition); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invalid MFA record"})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid MFA code"})
		return
	}

	err = db.DeleteMfaRecord(keyFilter, condition)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete MFA code"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message":  "User logged in successfully",
		"mfa_code": code,
	})

}
