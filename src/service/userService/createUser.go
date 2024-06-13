package userservice

import (
	"net/http"
	"strings"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/validator"
	"gorm.io/gorm"
)

func CreateNewUser(user model.UserJwtBody, newUser model.User) (u model.User, err error) {
    if err = user.ValidateRole(model.Admin); err != nil {
        return
    }

    newUser.Role = model.UserRole(strings.ToLower(string(newUser.Role)))
    println(newUser.Role)

    if err = validator.ValidateUser(&newUser); err != nil {
        return
    }

    _, err = databasehandler.GetUserByName(newUser.Username)
    if err != gorm.ErrRecordNotFound {
        err = cerror.New(http.StatusConflict, "OBJECT_ALREADY_EXISTS", "user already exists")
        return
    }

    dh := databasehandler.NewDatabaseHandler(user)
    defer func ()  {
        err = dh.CommitOrRollback(err)
    }()

    u, err = dh.CreateUser(newUser)
    u.Password = ""

    println(u.ID)

    return
}

func GetAllUsers(user model.UserJwtBody) (users []model.User, err error) {
    if err = user.ValidateRole(model.Admin); err != nil {
        return
    }

    users, err = databasehandler.GetAllUsers()
    return
}

func DeleteUser(user model.UserJwtBody, userId uint) (err error) {
    if err = user.ValidateRole(model.Admin); err != nil {
        return
    }

    dh := databasehandler.NewDatabaseHandler(user)
    defer func ()  {
        err = dh.CommitOrRollback(err)
    }()

    err = dh.DeleteUser(userId)
    return
}
