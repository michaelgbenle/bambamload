package models

import "time"

type User struct {
	Model
	Name         string `json:"name" gorm:"type:varchar(100)"`
	Email        string `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	PhoneNumber  string `json:"phone_number" gorm:"type:varchar(20);uniqueIndex"`
	BusinessName string `json:"business_name" gorm:"type:varchar(100)"`

	KycStatus            string `json:"kyc_status" gorm:"type:varchar(50);default:'business_profile'"` //business_profile,identity_verification,in_review,approved
	DocumentUploadStatus string `json:"document_upload_status" gorm:"type:varchar(20);default:'pending'"`

	//supplier details
	AccountType         string `json:"account_type" gorm:"type:varchar(100)"`
	BusinessDescription string `json:"business_description" gorm:"type:varchar(2550)"`
	YearFounded         string `json:"year_founded" gorm:"type:varchar(20)"`
	WebsiteUrl          string `json:"website_url" gorm:"type:varchar(255)"`
	LinkedInProfile     string `json:"linked_in_profile" gorm:"type:varchar(255)"`
	Country             string `json:"country" gorm:"type:varchar(100)"`
	State               string `json:"state" gorm:"type:varchar(100)"`
	Address             string `json:"address" gorm:"type:varchar(255)"`
	RegionsServed       string `json:"regions_served" gorm:"type:varchar(500)"`

	CacCertificate         string `json:"cac_certificate" gorm:"type:varchar(500)"`
	CacCertificateApproved bool   `json:"cac_certificate_approved" gorm:"default:false"`
	CacCertificateComment  string `json:"cac_certificate_comment" gorm:"type:varchar(200)"`

	ValidPersonalID         string `json:"valid_personal_id" gorm:"type:varchar(500)"`
	ValidPersonalIDApproved bool   `json:"valid_personal_id_approved" gorm:"default:false"`
	ValidPersonalIDComment  string `json:"valid_personal_id_comment" gorm:"type:varchar(200)"`

	UtilityBill         string `json:"utility_bill" gorm:"type:varchar(500)"`
	UtilityBillApproved bool   `json:"utility_bill_approved" gorm:"default:false"`
	UtilityBillComment  string `json:"utility_bill_comment" gorm:"type:varchar(200)"`

	TinDocument         string `json:"tin_document" gorm:"type:varchar(500)"`
	TinDocumentApproved bool   `json:"tin_document_approved" gorm:"default:false"`
	TinDocumentComment  string `json:"tin_document_comment" gorm:"type:varchar(200)"`

	SupplierRejectReason string  `json:"supplier_reject_reason" gorm:"type:varchar(200)"`
	CommissionRate       float32 `json:"commission_rate" gorm:"type:decimal(10,2)"`

	Employer  string `json:"employer" gorm:"type:varchar(100)"`
	Bvn       string `json:"bvn" gorm:"type:varchar(15);uniqueIndex"`
	Nin       string `json:"nin" gorm:"type:varchar(15);uniqueIndex"`
	Reference string `json:"reference" gorm:"type:varchar(255)"`

	Password      string    `json:"-" gorm:"type:varchar(255);not null"`
	Status        string    `json:"status" gorm:"type:varchar(50)"`
	Role          string    `json:"role" gorm:"type:varchar(50)"`
	IsBlocked     bool      `json:"is_blocked" gorm:"default:false"`
	IsActive      bool      `json:"is_active" gorm:"default:false"`
	LastLoginTime time.Time `json:"last_login_time" gorm:"type:timestamp"`
}

type SupplierRegisterRequest struct {
	Password  string `json:"password"`
	Reference string `json:"reference"`
}

type BuyerRegisterRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type VerifyForgotPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type VerifyRegistrationOtpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResendOtpRequest struct {
	Email  string `json:"email"`
	Action string `json:"action"`
}

type InviteSupplier struct {
	BusinessName      string `json:"business_name"`
	ContactPerson     string `json:"contact_person"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
	InvitationMessage string `json:"invitation_message"`
}

type ApproveOrRejectSupplierRequest struct {
	Action      string `json:"action"`
	Comment     string `json:"comment"`
	SupplierID  string `json:"supplier_id"`
	DocumentKey string `json:"document_key"`
}

// UpdateUserRequest represents the payload for user update.
// swagger:model
type UpdateUserRequest struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	DateOfBirth    string `json:"date_of_birth"`
	Nationality    string `json:"nationality"`
	Occupation     string `json:"occupation"`
	Employer       string `json:"employer"`
	Bvn            string `json:"bvn"`
	Nin            string `json:"nin"`
	Email          string `json:"email"`
	FaceImage      string `json:"face_image"`
	PhoneNumber    string `json:"phone_number"`
	ReferralCode   string `json:"referral_code"`
	TransactionPin string `json:"transaction_pin"`
}

type SignInRes struct {
	AccessToken       string    `json:"access_token"`
	AccessTokenExpiry time.Time `json:"access_token_expiry"`
	User              User      `json:"user"`
}

// PasswordChangeRequest represents a password change request
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// OTPVerificationRequest represents an OTP verification request
type OTPVerificationRequest struct {
	OTP string `json:"otp"`
}

type PinChangeRequest struct {
	CurrentPin string `json:"current_pin"`
	NewPin     string `json:"new_pin"`
	ConfirmPin string `json:"confirm_pin"`
}

type SupplierStats struct {
	TotalSuppliers int64
	Invited        int64
	Registering    int64
	Pending        int64
	Verified       int64
}
type SupplierAgg struct {
	Status string
	Count  int64
}

type EditSupplierCommissionRate struct {
	SupplierID string  `json:"supplier_id"`
	Rate       float32 `json:"rate"`
}
