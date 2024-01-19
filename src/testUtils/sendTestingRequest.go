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

type TestEnv struct {
    Router *gin.Engine
}

func InitRouter(reset bool) TestEnv {
    testEnv := TestEnv{}
    testEnv.Router = router.InitRouter()
    gin.SetMode(gin.TestMode)
    database.SetupDatabaseConnection()
    database.InitDatabaseStructure()

    if reset {
        debug.TruncateDatabase()
    }

    return testEnv
}

// deprecated
func SendTestingRequest(t *testing.T, req *http.Request, prepFunc ...func()) *httptest.ResponseRecorder  {
    env := InitRouter(true)

    for _, prep := range prepFunc {
        prep()
    }

    w := httptest.NewRecorder()
    env.Router.ServeHTTP(w, req)

    return w
}

func (env *TestEnv) SendTestingRequestSuccess(t *testing.T, req *http.Request, prepFunc func(), expectedHttpCode int, validateBody bool) api.SuccessResponse {
    prepFunc()

    w := httptest.NewRecorder()
    env.Router.ServeHTTP(w, req)

    assert.Equal(t, expectedHttpCode, w.Code)

    if(validateBody) {
        res := api.SuccessResponse{}
        err := json.Unmarshal(w.Body.Bytes(), &res)
        assert.Nilf(t, err, "Could now Unmarshal json: %s", w.Body.String())

        assert.Equal(t, true, res.Success, w.Body.String())

        return res
    }

    return api.SuccessResponse{}
}


func (env *TestEnv) SendTestingRequestError(t *testing.T, req *http.Request, prepFunc func(), expectedHttpCode int, expectedErrorType string) api.ErrorResponse {
    prepFunc()

    w := httptest.NewRecorder()
    env.Router.ServeHTTP(w, req)

    assert.Equal(t, expectedHttpCode, w.Code)

    res := api.ErrorResponse{}
    err := json.Unmarshal(w.Body.Bytes(), &res)
    assert.Nilf(t, err, "Could now Unmarshal json: %s", w.Body.String())

    assert.Equal(t, false, res.Success, res.Success)
    assert.Equal(t, expectedErrorType, res.Type, res.Type)

    return res
}


func ParseAndValidateResponse(t *testing.T, w *httptest.ResponseRecorder) api.SuccessResponse {
    res := api.SuccessResponse{}
    err := json.Unmarshal(w.Body.Bytes(), &res)

    assert.Nilf(t, err, "Could now Unmarshal json: %s", w.Body.String())

    return res
}
