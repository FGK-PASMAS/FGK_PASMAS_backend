package testutils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/database/debug"
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func SendTestingRequest(t *testing.T, req *http.Request, prepFunc ...func()) *httptest.ResponseRecorder  {
    r := router.InitRouter()
    gin.SetMode(gin.TestMode)
    database.SetupDatabaseConnection()
    database.InitDatabaseStructure()
    debug.TruncateDatabase()

    for _, prep := range prepFunc {
        prep()
    }

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    return w
}


func ParseAndValidateResponse(t *testing.T, w *httptest.ResponseRecorder) api.SuccessResponse {
    res := api.SuccessResponse{}
    err := json.Unmarshal(w.Body.Bytes(), &res)

    assert.Nilf(t, err, "Could now Unmarshal json: %s", w.Body.String())

    return res
}
