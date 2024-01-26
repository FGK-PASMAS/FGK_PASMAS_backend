package plane

import (
	"net/http"

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

    includes, err := pasmasservice.ParsePlaneInclude(c)
    if err == nil {
        planes, err = pasmasservice.GetPlanes(includes, nil)
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
