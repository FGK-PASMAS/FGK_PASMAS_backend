package pasmasservice

import (
	"time"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// This is vor debug only
var (
    key = "javainuse-secret-key"
    iss = "this_is_a_iss_key"
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
        claims := jwt.MapClaims{
            "username": user.Name,
            "role": user.Role,
            "iss": iss,
        }
        t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

        token, signingErr := t.SignedString([]byte(key))
        if signingErr == nil {
            return token, nil
        }
    }

    waitTime := startTime.Add(2000* time.Millisecond).Sub(time.Now().UTC())
    println("Waiting for: ", waitTime)
    time.Sleep(waitTime)

    return 
}

func ValidateJwt(token string) (user model.UserJwtBody, err error){
    claims := jwt.MapClaims { }
    tok, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
        return []byte(key), nil
    })

    if err != nil || tok.Valid == false {
        err = cerror.ErrInvalidCredentials
        return
    }


    return
}

func checkPasswordHash(password string, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

    return err == nil
}
