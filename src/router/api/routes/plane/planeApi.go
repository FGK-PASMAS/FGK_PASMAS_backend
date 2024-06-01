package plane

import (
	"net/http"
	"strconv"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/MetaEMK/FGK_PASMAS_backend/service/planeService"
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
		planes, err = planeService.GetPlanes(includes, filters)
	} else {
		if incErr != nil {
			err = incErr
		} else {
			err = filtErr
		}
	}

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: planes,
		}

		httpCode = http.StatusOK
	}

	c.JSON(httpCode, response)
}

func getPlaneById(c *gin.Context) {
	var response interface{}
	var httpCode int

	var planes model.Plane
	var err error

	includes, err := databasehandler.ParsePlaneInclude(c)

	if err == nil {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err == nil {
			planes, err = planeService.GetPlaneById(uint(id), includes)
		}
	}

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: planes,
		}

		httpCode = http.StatusOK
	}

	c.JSON(httpCode, response)
}

func updatePlane(c *gin.Context) {
	var response interface{}
	var httpCode int

	var plane model.Plane
	var id uint64
	var err error

	user := c.Keys["user"].(model.UserJwtBody)

	idStr := c.Param("id")
	id, err = strconv.ParseUint(idStr, 10, 64)

	if err == nil {
		var updateData databasehandler.PartialUpdatePlaneStruct
		err = c.ShouldBind(&updateData)

		if err == nil {
			plane, err = planeService.UpdatePlane(user, uint(id), updateData)
		}
	}

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: plane,
		}

		httpCode = http.StatusOK
	}

	c.JSON(httpCode, response)
}
