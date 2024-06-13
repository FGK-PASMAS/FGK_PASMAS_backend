package passenger

import (
	"net/http"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/config"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/passengerService"
	"github.com/gin-gonic/gin"
)

var log = logging.NewLogger("PassengerAPI", config.GetGlobalLogLevel())

func getPassengers(c *gin.Context) {
	passengers, err := passengerService.GetPassengers()

	if err != nil {
		apiErr := cerror.InterpretError(err)
		c.JSON(apiErr.HttpCode, apiErr)
	} else {
		c.JSON(http.StatusOK, api.SuccessResponse{Success: true, Response: passengers})
	}

}
