package models

type Misc_Details struct {
	Id          string `json:"id,omitempty"`
	Type_Code   string `json:"type_code"`
	Description string `json:"description"`
	Key_1       string `json:"key_1"`
	Val_1       string `json:"val_1"`
	Key_2       string `json:"key_2"`
	Val_2       string `json:"val_2"`
}

type Mutual_Fund struct {
	Id                   string `json:"id,omitempty"`
	Fund_Name            string `json:"fund_name"`
	Folio_Number         string `json:"folio_number"`
	User_Id              string `json:"user_id"`
	Email_Id             string `json:"email_id"`
	Mobile_Number        string `json:"mobile_number"`
	Login_Password       string `json:"login_password"`
	Transaction_Password string `json:"transaction_password"`
	Mpin                 string `json:"mpin"`
}

type Email struct {
	Id       string `json:"id,omitempty"`
	Entity   string `json:"entity"`
	Email_Id string `json:"email_id"`
	Password string `json:"password"`
}

type IRCTC_Details struct {
	Id        string `json:"id,omitempty"`
	User_Name string `json:"user_name"`
	Email_Id  string `json:"email_id"`
	Password  string `json:"password"`
}

type Passport_Details struct {
	Id              string `json:"id,omitempty"`
	Passport_Number string `json:"passport_number"`
	Issuer_Country  string `json:"issuer_country"`
	Issue_Date      string `json:"issue_date"`
	Expiry_Date     string `json:"expiry_date"`
	User_Id         string `json:"user_id"`
	Password        string `json:"password"`
}

type Aadhaar_Details struct {
	Id             string `json:"id,omitempty"`
	Aadhaar_Number string `json:"aadhaar_number"`
	Issue_Date     string `json:"issue_date"`
}

type Pran_Details struct {
	Id          string `json:"id,omitempty"`
	Pran_Number string `json:"pran_number"`
	Password    string `json:"password"`
}

type Pan_Details struct {
	Id         string `json:"id,omitempty"`
	Pan_Number string `json:"pan_number"`
	Issue_Date string `json:"issue_date"`
	Password   string `json:"password"`
}

type Bank_Details struct {
	Id                     string `json:"id,omitempty"`
	User_Id                string `json:"user_id"`
	Customer_Id            string `json:"customer_id"`
	Bank_Name              string `json:"bank_name"`
	Login_Password         string `json:"login_password"`
	Transaction_Password   string `json:"transaction_password"`
	Mobile_Login_Pin       string `json:"mobile_login_pin"`
	Mobile_Transaction_Pin string `json:"mobile_transaction_pin"`
}

type Address struct {
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Postal_Code string `json:"postal_code"`
	Country     string `json:"country"`
}

type User struct {
	User_Name        string             `json:"user_name"`
	Is_Admin         string             `json:"is_admin"`
	First_name       string             `json:"first_name"`
	Last_name        string             `json:"last_name"`
	Email            string             `json:"email"`
	Mobile_Number    string             `json:"mobile_number"`
	Password         string             `json:"password"`
	Date_of_birth    string             `json:"date_of_birth"`
	Address          Address            `json:"address"`
	Pan_Details      []Pan_Details      `json:"pan_details,omitempty"`
	Aadhaar_Details  []Aadhaar_Details  `json:"aadhaar_details,omitempty"`
	Pran_Details     []Pran_Details     `json:"pran_details,omitempty"`
	Passport_Details []Passport_Details `json:"passport_details,omitempty"`
	IRCTC_Details    []IRCTC_Details    `json:"irctc_details,omitempty"`
	Emails           []Email            `json:"emails,omitempty"`
	Mutual_Funds     []Mutual_Fund      `json:"mutual_funds,omitempty"`
	Misc_Details     []Misc_Details     `json:"misc_details,omitempty"`
	Bank_Details     []Bank_Details     `json:"bank_details,omitempty"`
}

type PasswordsInput struct {
	Password             string `json:"existing_password"`
	New_Password         string `json:"password"`
	Confirm_New_Password string `json:"confirm_password"`
}

type MFA_Code struct {
	User_Name string `json:"user_name"`
	Email_Id  string `json:"email_id"`
	Code      string `json:"code"`
}
