package apis

import (
	"go-mongo-project/config"
	"go-mongo-project/models"
	"go-mongo-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

/** User Related APIs **/

func RegisterUser(context *gin.Context) {

	user := models.User{}

	// Bind JSON to user struct
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build user data map (excluding user_name and password which are handled separately)
	userData := bson.M{
		config.USER_EMAIL:         user.Email,
		config.USER_MOBILE_NUMBER: user.Mobile_Number,
	}

	// Encrypt user-level fields
	encryptedUserData, err := utils.EncryptStructFields(userData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt user details"})
		return
	}

	// Update user object with encrypted data
	user.Email = encryptedUserData[config.USER_EMAIL].(string)
	user.Mobile_Number = encryptedUserData[config.USER_MOBILE_NUMBER].(string)

	// Build and encrypt address data
	if user.Address != (models.Address{}) {
		addressData := bson.M{
			config.ADDRESS_STREET:      user.Address.Street,
			config.ADDRESS_CITY:        user.Address.City,
			config.ADDRESS_STATE:       user.Address.State,
			config.ADDRESS_POSTAL_CODE: user.Address.Postal_Code,
			config.ADDRESS_COUNTRY:     user.Address.Country,
		}

		encryptedAddressData, err := utils.EncryptStructFields(addressData)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt address details"})
			return
		}

		// Update address with encrypted data
		user.Address.Street = encryptedAddressData[config.ADDRESS_STREET].(string)
		user.Address.City = encryptedAddressData[config.ADDRESS_CITY].(string)
		user.Address.State = encryptedAddressData[config.ADDRESS_STATE].(string)
		user.Address.Postal_Code = encryptedAddressData[config.ADDRESS_POSTAL_CODE].(string)
		user.Address.Country = encryptedAddressData[config.ADDRESS_COUNTRY].(string)
	}

	registeredUser := user

	// Save the user (password is already hashed in AddDocument method)
	if err := user.AddDocument(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResultSet, err := registeredUser.RetrieveUserData()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(userResultSet.Email, userResultSet.User_Name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	userResultSet.Password = ""

	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"token":   token,
		"result":  userResultSet,
	})

}

func Login(context *gin.Context) {
	// Login logic here
	loggedInUser := models.User{}

	if err := context.ShouldBindJSON(&loggedInUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userResultSet, err := loggedInUser.RetrieveUserData()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken(userResultSet.Email, userResultSet.User_Name)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	userResultSet.Password = ""

	context.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"token":   token,
		"result":  userResultSet,
	})

}

func UpdateUserPassword(context *gin.Context) {
	user := models.User{}

	input := models.PasswordsInput{}

	// Bind JSON to user struct
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.User_Name = context.GetString(config.USER_NAME)

	err := user.PasswordValidations(input)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})

}

func UpdateUserDetails(context *gin.Context) {

	user := models.User{}

	// Bind JSON to user struct
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.User_Name = context.GetString(config.USER_NAME)

	keyFilter := bson.M{config.USER_NAME: user.User_Name}

	// Build user data map
	userData := bson.M{
		config.USER_EMAIL:         user.Email,
		config.USER_MOBILE_NUMBER: user.Mobile_Number,
	}

	// Encrypt user-level fields (exclude user_name as it's used for filtering)
	encryptedUserData, err := utils.EncryptStructFields(userData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt user details"})
		return
	}

	// Build address data map
	addressData := bson.M{
		config.ADDRESS_STREET:      user.Address.Street,
		config.ADDRESS_CITY:        user.Address.City,
		config.ADDRESS_STATE:       user.Address.State,
		config.ADDRESS_POSTAL_CODE: user.Address.Postal_Code,
		config.ADDRESS_COUNTRY:     user.Address.Country,
	}

	// Encrypt address fields
	encryptedAddressData, err := utils.EncryptStructFields(addressData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt address details"})
		return
	}

	// Combine into update document
	updatedData := bson.M{
		"$set": bson.M{
			config.USER_EMAIL:         encryptedUserData[config.USER_EMAIL],
			config.USER_MOBILE_NUMBER: encryptedUserData[config.USER_MOBILE_NUMBER],
			"address":                 encryptedAddressData,
		},
	}

	err = user.UpdateRecord(keyFilter, updatedData)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Respond with success
	context.JSON(http.StatusCreated, gin.H{
		"message": "User Details updated in the document successfully.",
	})

}

func AllUsersData(context *gin.Context) {

	user := models.User{}
	user.User_Name = context.GetString(config.USER_NAME)

	resultSet, err := user.GetAllUsersData()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Data retrieved successfully",
		"result":  resultSet,
	})

}

/** End User Related APIs **/
