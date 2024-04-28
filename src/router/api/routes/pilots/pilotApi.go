package pilots

import (
	"net/http"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/pilotService"
	"github.com/gin-gonic/gin"
)

func getPilots(c *gin.Context) {
    var response interface{}
    var httpCode int

    var pilots *[]model.Pilot
    var err error

    include, err := pilotService.ParsePilotInclude(c)
    if err == nil {
        filter, err := pilotService.ParsePilotFilter(c)
        if err == nil {
            pilots, err = pilotService.GetPilots(include, filter)
        }
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        response = res.ErrorResponse
        httpCode = res.HttpCode
    } else {
        response = api.SuccessResponse{
            Success: true,
            Response: pilots,
        }
        httpCode = http.StatusOK
    }

    c.JSON(httpCode, response)
}
