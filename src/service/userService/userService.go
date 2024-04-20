package pasmasservice

import (
	"time"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateJwtForUser(username string, password string) (token string, err error) {
    startTime := time.Now().UTC()
    var user model.User

    if username == "" || password == "" {
        err = cerror.ErrEmptyCredentials
    } else {
        user, err = databasehandler.GetUserByName(username)
        if err != nil {
            err = cerror.ErrInvalidCredentials
        }
    }

    if err == nil && checkPasswordHash(password, user.Password) {
        claims := user.ToJwtClaims()
        claims["iss"] = config.JwtIssuer
        t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        token, signingErr := t.SignedString([]byte(config.JwtEncodingKey))
        if signingErr == nil {
            return token, nil
        }
    }
    
    err = cerror.ErrInvalidCredentials

    waitTime := startTime.Add(2000* time.Millisecond).Sub(time.Now().UTC()).Abs()
    println("Waiting for: ", waitTime)
    time.Sleep(waitTime)

    return 
}

func checkPasswordHash(password string, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

    return err == nil
}
