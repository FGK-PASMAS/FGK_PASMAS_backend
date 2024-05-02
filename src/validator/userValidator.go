package validator

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


func ValidateUser(user *model.User) error {
    if user.Username == "" {
        return cerror.ErrInvalidUsername
    }

    if user.Password == "" {
        return cerror.ErrInvalidPassword
    }

    if user.Role != model.Admin && user.Role != model.Vendor && user.Role != model.ReadOnly {
        return cerror.ErrInvalidRole
    }

    return nil
}
