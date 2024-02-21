package middleware

import (
	"net/http"
	"strings"

	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/gin-gonic/gin"
)

func ValidateJwt(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
    token := strings.Replace(auth, "Bearer ", "", 1)
    user, err := pasmasservice.ValidateJwt(token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, nil)
        c.Abort()
    }

    println(user.Username)

    c.Next()
}
