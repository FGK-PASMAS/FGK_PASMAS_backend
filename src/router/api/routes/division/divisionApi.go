package division

import (
	"net/http"

	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/divisionService"
	"github.com/gin-gonic/gin"
)

func getDivisions(c *gin.Context) {
    divisions, err := divisionService.GetDivisions()

    if err != nil {
        apiErr := api.GetErrorResponse(err)
        apiErr.ErrorResponse.Message = err.Error()
        c.JSON(apiErr.HttpCode, apiErr.ErrorResponse)
    } else {
        c.JSON(http.StatusOK, api.SuccessResponse{Success: true, Response: divisions})
    }
}
