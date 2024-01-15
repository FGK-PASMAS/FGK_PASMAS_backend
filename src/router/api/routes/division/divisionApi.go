package division

import (
	"net/http"

	divisionhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/divisionHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/gin-gonic/gin"
)

func getDivision(c *gin.Context) {
    var statusCode int
    var response interface{}

    divisions, err := divisionhandler.GetDivision()

    if err != nil {
        statusCode = http.StatusInternalServerError
        response = api.ErrorResponse{Success: false, ErrorBody: err}
    } else {
        statusCode = http.StatusOK
        response = api.SuccessResponse{Success: true, Response: divisions}
    }

    c.JSON(statusCode, response)
}
