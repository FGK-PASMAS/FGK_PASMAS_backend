package division

import (
	"net/http"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/divisionService"
	"github.com/gin-gonic/gin"
)

func getDivisions(c *gin.Context) {
    divisions, err := divisionService.GetDivisions()

    if err != nil {
        apiErr := cerror.InterpretError(err)
        c.JSON(apiErr.HttpCode, apiErr)
    } else {
        c.JSON(http.StatusOK, api.SuccessResponse{Success: true, Response: divisions})
    }
}
