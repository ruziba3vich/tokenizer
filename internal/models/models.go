package models

type OneTimeLink struct {
	ID   uint   `gorm:"primaryKey"`
	Key  string `gorm:"uniqueIndex;not null"`
	Used bool   `gorm:"default:false"`
}

// RegisterPayload represents the user registration input
type RegisterPayload struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
	Phone     string `json:"phone,omitempty" example:"+998901234567"`
	Username  string `json:"username" example:"johndoe"`
	Password  string `json:"password" example:"securepassword123"`
	Key       string `json:"key" example:"abc123"`
}

// LoginPayload represents login input
type LoginPayload struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"securepassword123"`
}

// User represents the returned user
type User struct {
	ID        uint   `json:"id" example:"1"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"john@example.com"`
	Phone     string `json:"phone,omitempty" example:"+998901234567"`
	Username  string `json:"username" example:"johndoe"`
	Password  string `json:"password" example:"password"`
}
