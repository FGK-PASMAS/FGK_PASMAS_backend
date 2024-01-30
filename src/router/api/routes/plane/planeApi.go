package plane

import (
	"net/http"
	"strconv"

	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	pasmasservice "github.com/MetaEMK/FGK_PASMAS_backend/service/pasmasService"
	"github.com/gin-gonic/gin"
)

func getPlanes(c *gin.Context) {
    var response interface{}
    var httpCode int

    var planes *[]model.Plane
    var err error

    includes, incErr := pasmasservice.ParsePlaneInclude(c)
    filters, filtErr := pasmasservice.ParsePlaneFilter(c)

    if incErr == nil && filtErr == nil {
        planes, err = pasmasservice.GetPlanes(includes, filters)
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

func updatePlane(c *gin.Context) {
    var response interface{}
    var httpCode int
    var plane *model.Plane
    var err error

    planeIdStr := c.Param("id")
    planeId, err := strconv.ParseUint(planeIdStr, 10, 64)

    if err == nil {
        var body struct{
            PrefPilotId uint
        }

        err = c.ShouldBind(&body)

        if err == nil {
            plane, err = pasmasservice.UpdatePrefPilot(uint(planeId), body.PrefPilotId)
        }

    }


    if err != nil {
        res := api.GetErrorResponse(err)
        response = res.ErrorResponse
        httpCode = res.HttpCode
    } else {
        response = api.SuccessResponse{
            Success: true,
            Response: plane,
        }

        httpCode = http.StatusOK
    }

    c.JSON(httpCode, response)
}
