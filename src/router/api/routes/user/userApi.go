package user

import (
	"net/http"

	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	userservice "github.com/MetaEMK/FGK_PASMAS_backend/service/userService"
	"github.com/gin-gonic/gin"
)

func ValidateUser(c *gin.Context) {
    var httpCode int
    var response interface{}
    var err error

    var token string

    username, password, ok := c.Request.BasicAuth()

    if ok {
        token, err = userservice.GenerateJwtForUser(username, password)
        println("Token: ", token)
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        httpCode = res.HttpCode
        response = res.ErrorResponse
    } else {
        response = api.SuccessResponse {
            Success: true,
            Response: token,
        }
        httpCode = http.StatusOK
    }
    c.JSON(httpCode, response)
}
