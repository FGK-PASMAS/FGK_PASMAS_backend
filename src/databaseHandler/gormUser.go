package databasehandler

import (

    "github.com/MetaEMK/FGK_PASMAS_backend/model"
    "golang.org/x/crypto/bcrypt"
)


func initUser() {
    Db.AutoMigrate(&model.User{})

    dh := NewDatabaseHandler()
    defer dh.CommitOrRollback(nil)

    admin := model.User{
        Name: "admin",
        Role: model.Admin,
        Password: "admin123",
    }

    vendor := model.User{
        Name: "vendor",
        Role: model.Vendor,
        Password: "vendor123",
    }

    readOnly := model.User{
        Name: "readOnly",
        Role: model.ReadOnly,
        Password: "readOnly123",
    }

    dh.CreateUser(admin)
    dh.CreateUser(vendor)
    dh.CreateUser(readOnly)
}


func (dh *DatabaseHandler) CreateUser(user model.User) (newUser model.User, err error) {
    passwordHash, err := hashPassword(user.Password)
    if err != nil {
        return
    }
    user.Password = passwordHash

    err = dh.Db.Create(&user).Error
    return
}

func hashPassword(password string) (hash string, err error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

    hash = string(bytes)
    return
}

func checkPasswordHash(password string, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

    return err == nil
}
