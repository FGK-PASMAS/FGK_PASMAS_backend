package model

import "gorm.io/gorm"


type User struct {
    gorm.Model

    // Name of the user. This includes first name and last name or a username.
    Name            string             `gorm:"not null"`

    // Password of the user. This is a hashed password.
    Password        string             `gorm:"not null"`

    // Role of the user. 
    Role            UserRole           `gorm:"not null"`
}

type UserJwtBody struct {
    Username        string
    Role            UserRole
}

type UserRole string

const (
    Admin       UserRole = "admin"
    Vendor      UserRole = "user"
    ReadOnly    UserRole = "read-only"
)
