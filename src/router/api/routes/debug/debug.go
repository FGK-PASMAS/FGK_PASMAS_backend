package debug

import (
	"net/http"

	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
    response := api.SuccessResponse{
        Success: true,
        Response: "pong",
    }

    c.JSON(http.StatusOK, response)
}

func healthCheck(c *gin.Context) {

    c.JSON(http.StatusOK, api.SuccessResponse{Success: true, Response: "TODO"})
}

func resetDatabase(c *gin.Context) {
    err := pasmasservice.TruncateData()

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        c.JSON(http.StatusOK, api.SuccessResponse{Success: true})
    }
}
