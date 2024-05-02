package model

import (

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)


type User struct {
    gorm.Model

    // Username of the user. This includes first name and last name or a username.
    Username            string        `gorm:"not null"`

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
    Vendor      UserRole = "vendor"
    ReadOnly    UserRole = "read-only"
)

var permissionList = map[UserRole]int {
    Admin:      3,
    Vendor:     2,
    ReadOnly:   1,
}

func (user *UserJwtBody) ValidateRole(neededRole UserRole) error {
    if permissionList[neededRole] <= permissionList[user.Role] {
        return nil
    }

    return cerror.ErrForbidden
}

func (user *User) ToJwtClaims() jwt.MapClaims {
    return jwt.MapClaims{
        "Role": user.Role,
        "Username": user.Username,
    }
}

func (u *User) SetTimesToUTC() {
    u.CreatedAt = u.CreatedAt.UTC()
    u.UpdatedAt = u.UpdatedAt.UTC()
}

func ClaimsToUserJwtBody(claims jwt.MapClaims) (UserJwtBody, error) {
    var body UserJwtBody
    roleClaim, ok := claims["Role"].(string)
    if !ok {
        return body, cerror.ErrInvalidRole
    }
    role, err := convertStringToRole(roleClaim)
    username, ok := claims["Username"].(string)

    if ok && err == nil && len(username) > 0{
        body = UserJwtBody{
            Username: username,
            Role: role,
        }
    }

    return body, err
}

// Convert a string to a UserRole
func convertStringToRole(str string) (role UserRole, err error) {
    switch str {
    case string(Admin):
        role = Admin
    case string(Vendor):
        role = Vendor
    case string(ReadOnly):
        role = ReadOnly
    default:
        err = cerror.ErrInvalidRole
    }

    return 
}
