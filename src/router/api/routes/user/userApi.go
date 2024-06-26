package user

import (
	"net/http"
	"strconv"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	userservice "github.com/MetaEMK/FGK_PASMAS_backend/service/userService"
	"github.com/gin-gonic/gin"
)

func getAllUsers(c *gin.Context) {
	var httpCode int
	var response interface{}
	var err error

	user := c.Keys["user"].(model.UserJwtBody)

	users, err := userservice.GetAllUsers(user)

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: users,
		}
		httpCode = http.StatusOK
	}

	c.JSON(httpCode, response)
}

func createNewUser(c *gin.Context) {
	var httpCode int
	var response interface{}
	var err error

	user := c.Keys["user"].(model.UserJwtBody)

	newUser := model.User{}
	err = c.ShouldBind(&newUser)

	newUser, err = userservice.CreateNewUser(user, newUser)

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success:  true,
			Response: newUser,
		}
	}

	c.JSON(httpCode, response)
}

func deleteUser(c *gin.Context) {
	var httpCode int
	var response interface{}
	var err error

	user := c.Keys["user"].(model.UserJwtBody)

	userIdStr := c.Param("id")
	userId, err := strconv.ParseUint(userIdStr, 10, 64)

	err = userservice.DeleteUser(user, uint(userId))

	if err != nil {
		e := cerror.InterpretError(err)
		httpCode = e.HttpCode
		response = e
	} else {
		response = api.SuccessResponse{
			Success: true,
		}
	}

	c.JSON(httpCode, response)
}
