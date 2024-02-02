package debug

import (
	"net/http"

	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
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
    err := databasehandler.ResetDatabase()

    if err != nil {
        res := api.GetErrorResponse(err)
        c.JSON(res.HttpCode, res.ErrorResponse)
    } else {
        c.JSON(http.StatusOK, api.SuccessResponse{Success: true})
    }
}
