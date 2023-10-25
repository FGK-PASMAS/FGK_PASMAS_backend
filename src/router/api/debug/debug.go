package debug

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
    c.String(http.StatusOK, "pong");
}
