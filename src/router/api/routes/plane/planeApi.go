package plane

import (
	"net/http"

	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	planeservice "github.com/MetaEMK/FGK_PASMAS_backend/service/planeService"
	"github.com/gin-gonic/gin"
)

func getPlanes(c *gin.Context) {
    var response interface{}
    var httpCode int

    var planes []model.Plane
    var err error

    includes, incErr := databasehandler.ParsePlaneInclude(c)
    filters, filtErr := databasehandler.ParsePlaneFilter(c)

    if incErr == nil && filtErr == nil {
        planes, err = planeservice.GetPlanes(includes, filters)
    } else {
        if incErr != nil {
            err = incErr
        } else {
            err = filtErr
        }
    }

    if err != nil {
        res := api.GetErrorResponse(err)
        response = res.ErrorResponse
        httpCode = res.HttpCode
    } else {
        response = api.SuccessResponse{
            Success: true,
            Response: planes,
        }

        httpCode = http.StatusOK
    }

    c.JSON(httpCode, response)
}
