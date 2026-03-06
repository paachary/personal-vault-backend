package routes

import (
	"go-mongo-project/apis"
	"go-mongo-project/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	// User Related
	server.POST("/register-user", apis.RegisterUser)
	server.POST("/login", apis.Login)
	server.POST("/request-mfa", apis.GenerateMFA)
	server.POST("/verify-mfa", apis.VerifyMFA)

	// Middleware

	authenicated := server.Group("/")
	authenicated.Use(middlewares.Authenticate)

	// PAN related
	authenicated.POST("/pan", apis.AddPANDetails)
	authenicated.PUT("/pan", apis.UpdatePANDetails)
	authenicated.DELETE("/pan", apis.DeletePANDetails)

	// Aadhaar related
	authenicated.POST("/aadhaar", apis.AddAadhaarDetails)
	authenicated.PUT("/aadhaar", apis.UpdateAadhaarDetails)
	authenicated.DELETE("/aadhaar", apis.DeleteAadhaarDetails)

	// PRAN related
	authenicated.POST("/pran", apis.AddPRANDetails)
	authenicated.PUT("/pran", apis.UpdatePRANDetails)
	authenicated.DELETE("/pran", apis.DeletePRANDetails)

	// Passport related
	authenicated.POST("/passport", apis.AddPassportDetails)
	authenicated.PUT("/passport", apis.UpdatePassportDetails)
	authenicated.DELETE("/passport", apis.DeletePassportDetails)

	// IRCTC related
	authenicated.POST("/irctc", apis.AddIRCTCDetails)
	authenicated.PUT("/irctc", apis.UpdateIRCTCDetails)
	authenicated.DELETE("/irctc", apis.DeleteIRCTCDetails)

	// Emails related
	authenicated.POST("/email", apis.AddEmailDetails)
	authenicated.PUT("/email", apis.UpdateEmailDetails)
	authenicated.DELETE("/email", apis.DeleteEmailDetails)

	// Mutual Funds Related
	authenicated.POST("/mf", apis.AddMFDetails)
	authenicated.PUT("/mf", apis.UpdateMFDetails)
	authenicated.DELETE("/mf", apis.DeleteMFDetails)

	// Bank Related
	authenicated.POST("/bank", apis.AddBankDetails)
	authenicated.PUT("/bank", apis.UpdateBankDetails)
	authenicated.DELETE("/bank", apis.DeleteBankDetails)

	// Misc Related
	authenicated.POST("/misc", apis.AddMiscDetails)
	authenicated.PUT("/misc", apis.UpdateMiscDetails)
	authenicated.DELETE("/misc", apis.DeleteMiscDetails)

	// Get Details
	authenicated.GET("/", apis.GetPersonalDetails)

	//User updates
	authenicated.PUT("/user-details", apis.UpdateUserDetails)
	authenicated.PUT("/user-password", apis.UpdateUserPassword)

	// All Users
	authenicated.GET("/all-users", apis.AllUsersData)

}
