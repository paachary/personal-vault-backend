package models

import (
	"errors"
	"fmt"
	"go-mongo-project/config"
	"go-mongo-project/db"
	"go-mongo-project/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("field '%s' has error: %s", e.Field, e.Message)
}

func _validateInputs(errorList *[]error, field any, fieldStr string) {

	if field == "" {
		err := &ValidationError{
			Field:   fieldStr,
			Message: "value is nil or empty",
		}
		*errorList = append(*errorList, err)
	}

}

func validatePasswordStrength(errorList *[]error, password string) {
	if len(password) < 8 {
		err := &ValidationError{
			Field:   "Password",
			Message: "must be at least 8 characters long",
		}
		*errorList = append(*errorList, err)
	}
}

func (u *User) AddDocument() error {

	var errorsList []error

	_validateInputs(&errorsList, u.User_Name, "User Name")
	_validateInputs(&errorsList, u.First_name, "First Name")
	_validateInputs(&errorsList, u.Last_name, "Last Name")
	_validateInputs(&errorsList, u.Date_of_birth, "Date of Birth")
	_validateInputs(&errorsList, u.Email, "Email")
	_validateInputs(&errorsList, u.Password, "Password")
	_validateInputs(&errorsList, u.Address.City, "City")
	_validateInputs(&errorsList, u.Address.State, "State")
	_validateInputs(&errorsList, u.Address.Country, "Country")
	_validateInputs(&errorsList, u.Address.Postal_Code, "postal code")
	_validateInputs(&errorsList, u.Address.Street, "Street")

	validatePasswordStrength(&errorsList, u.Password)

	if errorsList != nil {
		return errors.Join(errorsList...)
	}

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return errors.New("Failed to hash password: " + err.Error())
	}

	u.Password = hashedPassword

	_, err = db.InsertOne(u, config.USER_NAME, u.User_Name)

	if err != nil {

		return errors.New("Error occurred during adding the document : " + err.Error())
	}

	return nil

}

func (u *User) RetrieveUserData() (*User, error) {
	filter := bson.M{config.USER_NAME: u.User_Name}

	userResult, err := db.RetrieveSingleRecord[User](filter)
	if err != nil {
		return nil, errors.New("Invalid Credentials. Please resubmit.")
	}

	if err := userResult.DecryptCredentials(); err != nil {
		return nil, errors.New("Error decrypting credentials: " + err.Error())
	}

	valid := utils.CheckPasswordHash(u.Password, userResult.Password)

	if !valid {
		return nil, errors.New("Invalid Credentials. Please resubmit.")
	}

	return userResult, nil
}

func (u *User) PasswordValidations(inputs PasswordsInput) error {

	u.Password = inputs.Password

	_, err := u.RetrieveUserData()

	if err != nil {
		return err
	}

	if u.Password == inputs.New_Password {
		return errors.New("New Password cannot be same as the old Password. Please resubmit.")
	}

	if inputs.New_Password != inputs.Confirm_New_Password {
		return errors.New("New Passwords don't match with each other. Please resubmit.")
	}

	if len(inputs.New_Password) < 8 {
		return errors.New("New Password must be at least 8 characters long. Please resubmit.")
	}

	// Hash the password before saving
	hashedNewPassword, err := utils.HashPassword(inputs.New_Password)
	if err != nil {
		return errors.New("Failed to hash password: " + err.Error())
	}

	updatedData := bson.M{
		"$set": bson.M{
			config.USER_PASSWORD: hashedNewPassword,
		},
	}

	filter := bson.M{config.USER_NAME: u.User_Name}

	err = u.UpdateRecord(filter, updatedData)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetPersonalDetails() (*User, error) {

	filter := bson.M{config.USER_NAME: u.User_Name}

	userResult, err := db.RetrieveSingleRecord[User](filter)

	if err != nil {
		return nil, errors.New("Error retrieving personal data: " + err.Error())
	}

	if err := userResult.DecryptCredentials(); err != nil {
		return nil, errors.New("Error decrypting credentials: " + err.Error())
	}

	return userResult, nil
}

// DecryptCredentials decrypts all stored credential passwords and encrypted fields in the User document.
func (u *User) DecryptCredentials() error {
	dec := func(field *string) error {
		if *field == "" {
			return nil // Skip empty fields
		}
		val, err := utils.DecryptCredential(*field)
		if err != nil {
			return err
		}
		*field = val
		return nil
	}

	// Decrypt user-level fields
	if err := dec(&u.Email); err != nil {
		return fmt.Errorf("failed to decrypt email: %w", err)
	}
	if err := dec(&u.Mobile_Number); err != nil {
		return fmt.Errorf("failed to decrypt mobile number: %w", err)
	}

	// Decrypt address fields
	if err := dec(&u.Address.Street); err != nil {
		return fmt.Errorf("failed to decrypt address street: %w", err)
	}
	if err := dec(&u.Address.City); err != nil {
		return fmt.Errorf("failed to decrypt address city: %w", err)
	}
	if err := dec(&u.Address.State); err != nil {
		return fmt.Errorf("failed to decrypt address state: %w", err)
	}
	if err := dec(&u.Address.Postal_Code); err != nil {
		return fmt.Errorf("failed to decrypt address postal code: %w", err)
	}
	if err := dec(&u.Address.Country); err != nil {
		return fmt.Errorf("failed to decrypt address country: %w", err)
	}

	// Decrypt Aadhaar Details
	for i := range u.Aadhaar_Details {
		if err := dec(&u.Aadhaar_Details[i].Aadhaar_Number); err != nil {
			return fmt.Errorf("failed to decrypt aadhaar number: %w", err)
		}
		if err := dec(&u.Aadhaar_Details[i].Issue_Date); err != nil {
			return fmt.Errorf("failed to decrypt aadhaar issue date: %w", err)
		}
	}

	// Decrypt PAN Details
	for i := range u.Pan_Details {
		if err := dec(&u.Pan_Details[i].Pan_Number); err != nil {
			return fmt.Errorf("failed to decrypt pan number: %w", err)
		}
		if err := dec(&u.Pan_Details[i].Issue_Date); err != nil {
			return fmt.Errorf("failed to decrypt pan issue date: %w", err)
		}
		if err := dec(&u.Pan_Details[i].Password); err != nil {
			return fmt.Errorf("failed to decrypt pan password: %w", err)
		}
	}

	// Decrypt PRAN Details
	for i := range u.Pran_Details {
		if err := dec(&u.Pran_Details[i].Pran_Number); err != nil {
			return fmt.Errorf("failed to decrypt pran number: %w", err)
		}
		if err := dec(&u.Pran_Details[i].Password); err != nil {
			return fmt.Errorf("failed to decrypt pran password: %w", err)
		}
	}

	// Decrypt Passport Details
	for i := range u.Passport_Details {
		if err := dec(&u.Passport_Details[i].Passport_Number); err != nil {
			return fmt.Errorf("failed to decrypt passport number: %w", err)
		}
		if err := dec(&u.Passport_Details[i].Issue_Date); err != nil {
			return fmt.Errorf("failed to decrypt passport issue date: %w", err)
		}
		if err := dec(&u.Passport_Details[i].Expiry_Date); err != nil {
			return fmt.Errorf("failed to decrypt passport expiry date: %w", err)
		}
		if err := dec(&u.Passport_Details[i].Issuer_Country); err != nil {
			return fmt.Errorf("failed to decrypt passport issuer country: %w", err)
		}
		if err := dec(&u.Passport_Details[i].User_Id); err != nil {
			return fmt.Errorf("failed to decrypt passport user id: %w", err)
		}
		if err := dec(&u.Passport_Details[i].Password); err != nil {
			return fmt.Errorf("failed to decrypt passport password: %w", err)
		}
	}

	// Decrypt IRCTC Details
	for i := range u.IRCTC_Details {
		if err := dec(&u.IRCTC_Details[i].User_Name); err != nil {
			return fmt.Errorf("failed to decrypt irctc username: %w", err)
		}
		if err := dec(&u.IRCTC_Details[i].Email_Id); err != nil {
			return fmt.Errorf("failed to decrypt irctc email: %w", err)
		}
		if err := dec(&u.IRCTC_Details[i].Password); err != nil {
			return fmt.Errorf("failed to decrypt irctc password: %w", err)
		}
	}

	// Decrypt Email Details
	for i := range u.Emails {
		if err := dec(&u.Emails[i].Entity); err != nil {
			return fmt.Errorf("failed to decrypt email entity: %w", err)
		}
		if err := dec(&u.Emails[i].Email_Id); err != nil {
			return fmt.Errorf("failed to decrypt email id: %w", err)
		}
		if err := dec(&u.Emails[i].Password); err != nil {
			return fmt.Errorf("failed to decrypt email password: %w", err)
		}
	}

	// Decrypt Mutual Fund Details
	for i := range u.Mutual_Funds {
		if err := dec(&u.Mutual_Funds[i].Folio_Number); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund folio number: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].Fund_Name); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund name: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].User_Id); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund user id: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].Login_Password); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund login password: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].Transaction_Password); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund transaction password: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].Email_Id); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund email id: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].Mobile_Number); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund mobile number: %w", err)
		}
		if err := dec(&u.Mutual_Funds[i].Mpin); err != nil {
			return fmt.Errorf("failed to decrypt mutual fund Mpin: %w", err)
		}
	}

	// Decrypt Bank Details
	for i := range u.Bank_Details {
		if err := dec(&u.Bank_Details[i].Bank_Name); err != nil {
			return fmt.Errorf("failed to decrypt bank name: %w", err)
		}
		if err := dec(&u.Bank_Details[i].Customer_Id); err != nil {
			return fmt.Errorf("failed to decrypt bank customer id: %w", err)
		}
		if err := dec(&u.Bank_Details[i].User_Id); err != nil {
			return fmt.Errorf("failed to decrypt bank user id: %w", err)
		}
		if err := dec(&u.Bank_Details[i].Login_Password); err != nil {
			return fmt.Errorf("failed to decrypt bank login password: %w", err)
		}
		if err := dec(&u.Bank_Details[i].Transaction_Password); err != nil {
			return fmt.Errorf("failed to decrypt bank transaction password: %w", err)
		}
		if err := dec(&u.Bank_Details[i].Mobile_Login_Pin); err != nil {
			return fmt.Errorf("failed to decrypt bank mobile login pin: %w", err)
		}
		if err := dec(&u.Bank_Details[i].Mobile_Transaction_Pin); err != nil {
			return fmt.Errorf("failed to decrypt bank mobile transaction pin: %w", err)
		}
	}

	// Decrypt Misc Details
	for i := range u.Misc_Details {
		if err := dec(&u.Misc_Details[i].Type_Code); err != nil {
			return fmt.Errorf("failed to decrypt misc type code: %w", err)
		}
		if err := dec(&u.Misc_Details[i].Description); err != nil {
			return fmt.Errorf("failed to decrypt misc description: %w", err)
		}
		if err := dec(&u.Misc_Details[i].Key_1); err != nil {
			return fmt.Errorf("failed to decrypt misc key 1: %w", err)
		}
		if err := dec(&u.Misc_Details[i].Val_1); err != nil {
			return fmt.Errorf("failed to decrypt misc val 1: %w", err)
		}
		if err := dec(&u.Misc_Details[i].Key_2); err != nil {
			return fmt.Errorf("failed to decrypt misc key 2: %w", err)
		}
		if err := dec(&u.Misc_Details[i].Val_2); err != nil {
			return fmt.Errorf("failed to decrypt misc val 2: %w", err)
		}
	}

	return nil
}

func (u *User) AddRecord(keyFilter bson.M, attribute string, data bson.M, extraFilters map[string]any, multicheck bool) error {

	if multicheck {
		exists, err := db.FindRecordCount(keyFilter, attribute)

		if err != nil {
			return errors.New("Multiple records check failed while adding the attribute " + attribute + " in the doc. Please revisit your filter criteria.")
		}

		if exists {
			return errors.New("There is an existing record already for " + attribute + " in the document. Not proceeding futher.")
		}
	}

	err := db.AddOne(keyFilter, attribute, data, extraFilters)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (u *User) UpdateArrayRecord(keyFilter bson.M, arrayField string, matchCondition bson.M, data map[string]interface{}) error {

	err := db.UpdateOneForArray(keyFilter, arrayField, matchCondition, data)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (u *User) UpdateRecord(keyfilter bson.M, update bson.M) error {

	err := db.UpdateOne(keyfilter, update)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil

}

func (u *User) DeleteRecord(keyFilter bson.M, arrayField string, matchCondition bson.M) error {

	err := db.DeleteOne(keyFilter, arrayField, matchCondition)

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (u *User) GetAllUsersData() (*[]User, error) {

	keyFilter := bson.M{config.USER_NAME: u.User_Name}

	filter := bson.M{config.USER_NAME: bson.M{"$ne": u.User_Name}}

	userResult, err := db.GetAllUsersData[User](keyFilter, filter)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Decrypt credentials for all users
	for i := range *userResult {
		if err := (*userResult)[i].DecryptCredentials(); err != nil {
			return nil, errors.New("Error decrypting credentials for user " + (*userResult)[i].User_Name + ": " + err.Error())
		}
	}

	return userResult, nil

}
