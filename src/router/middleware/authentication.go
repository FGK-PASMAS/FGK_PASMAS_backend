package middleware

import (
	"net/http"
	"strings"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJwt(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
    token := strings.Replace(auth, "Bearer ", "", 1)
    user, err := validateToken(token)

    if err != nil {
        err := api.ErrorResponse {
            Success: false,
            Type: "Unauthorized",
            Message: "Invalid token",
        }

        c.JSON(http.StatusUnauthorized, err)
        c.Abort()
        return
    }

    c.Keys = make(map[string]interface{})
    c.Keys["user"] = user

    c.Next()
}

func validateToken(token string) (user model.UserJwtBody, err error){
    claims := jwt.MapClaims { }
    tok, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
        return []byte(config.JwtEncodingKey), nil
    })

    user, err = model.ClaimsToUserJwtBody(claims)

    if err != nil || tok.Valid == false {
        err = cerror.ErrInvalidCredentials
    }

    return
}
