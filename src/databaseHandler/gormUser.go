package databasehandler

import (
	"os"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func initUser() {
	Db.AutoMigrate(&model.User{})

	err := Db.Where("role = ?", "admin").First(&model.User{}).Error
	if err != nil {
        log.Debug("Admin user not found, creating one")

		if err == gorm.ErrRecordNotFound {
            pw := os.Getenv("ADMIN_PASSWORD")
            if pw == "" {
                panic("ADMIN_PASSWORD is not set")
            }

			dh := NewDatabaseHandler(model.UserJwtBody{})
			defer func()  {
                err = dh.CommitOrRollback(err)
                if err != nil {
                    panic("Admin user creation failed")
                }
			}()

			admin := model.User{
				Username: "admin",
				Role:     model.Admin,
				Password: pw,
			}

            _, err = dh.CreateUser(admin)
		}
	}
}

func GetAllUsers() (users []model.User, err error) {
	err = Db.Order("id ASC").Find(&users).Error

	for i := range users {
		users[i].SetTimesToUTC()
		users[i].Password = ""
	}

	return
}

func (dh *DatabaseHandler) CreateUser(user model.User) (newUser model.User, err error) {
	user.SetTimesToUTC()
	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		return
	}
	user.Password = passwordHash

	err = dh.Db.Create(&user).Error

	newUser = user
	newUser.Password = ""
	return
}

func GetUserByName(name string) (user model.User, err error) {
	err = Db.Model(&model.User{}).Where("username = ?", name).First(&user).Error
	user.SetTimesToUTC()

	return
}

func hashPassword(password string) (hash string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	hash = string(bytes)
	return
}

func (dh *DatabaseHandler) DeleteUser(userId uint) (err error) {

	var user model.User
	err = dh.Db.First(&user, userId).Error

	if err != nil {
		return
	}

	err = dh.Db.Delete(&model.User{}, userId).Error

	return
}
