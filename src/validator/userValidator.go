package validator

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


func ValidateUser(user *model.User) error {
    if user.Username == "" {
        return cerror.NewInvalidRequestBodyError("Username not valid")
    }

    if user.Password == "" {
        return cerror.NewInvalidRequestBodyError("Password not valid")
    }

    if user.Role != model.Admin && user.Role != model.Vendor && user.Role != model.ReadOnly {
        return cerror.NewInvalidRequestBodyError("Role not valid")
    }

    return nil
}
