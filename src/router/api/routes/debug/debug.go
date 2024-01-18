package debug

import (
	"net/http"

	"github.com/MetaEMK/FGK_PASMAS_backend/database/debug"
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

func truncate(c *gin.Context) {
    err := debug.TruncateDatabase()

    if err != nil {
        c.JSON(http.StatusInternalServerError, api.ErrorResponse{Success: false, Message: err.Error()})
    } else {
        response := struct{
            Truncate string     `json:"truncate"`
            Seeding string      `json:"seeding"`
        }{
            Truncate: "successfull",
            Seeding: "successfull",
        }
        c.JSON(http.StatusOK, api.SuccessResponse{Success: true, Response: response})
    }
}
