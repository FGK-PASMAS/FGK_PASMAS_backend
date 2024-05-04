package userservice

import (
	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func CreateNewUser(user model.UserJwtBody, newUser model.User) (err error) {
    if err = user.ValidateRole(model.Admin); err != nil {
        return
    }

    _, err = databasehandler.GetUserByName(newUser.Name)
    if err != cerror.ErrObjectNotFound {
        err = cerror.ErrUserAlreadyExists
        return
    }

    println("Lets Go")

    dh := databasehandler.NewDatabaseHandler()
    defer func ()  {
        err = dh.CommitOrRollback(err)
    }()

    _, err = dh.CreateUser(newUser)

    return
}
