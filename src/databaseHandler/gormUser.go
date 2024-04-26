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

    dh.createUserIfNotExists(admin)
    dh.createUserIfNotExists(vendor)
    dh.createUserIfNotExists(readOnly)
}

func (dh * DatabaseHandler) createUserIfNotExists(user model.User) {
    user.SetTimesToUTC()
    var userCount int64 = -1
    dh.Db.Model(&model.User{}).Where("name = ?", user.Name).Count(&userCount)

    if userCount == 0 {
        dh.CreateUser(user)
    }
}

func (dh *DatabaseHandler) CreateUser(user model.User) (newUser model.User, err error) {
    user.SetTimesToUTC()
    passwordHash, err := hashPassword(user.Password)
    if err != nil {
        return
    }
    user.Password = passwordHash

    err = dh.Db.Create(&user).Error
    return
}

func GetUserByName(name string) (user model.User, err error) {
    err = Db.Model(&model.User{}).Where("name = ?", name).First(&user).Error
    user.SetTimesToUTC()

    return 
}

func hashPassword(password string) (hash string, err error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

    hash = string(bytes)
    return
}

