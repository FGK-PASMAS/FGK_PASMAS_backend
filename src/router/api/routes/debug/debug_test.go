package debug_test

import (
	"net/http"
	"testing"

	"github.com/MetaEMK/FGK_PASMAS_backend/router/api"
	testutils "github.com/MetaEMK/FGK_PASMAS_backend/testUtils"
	"github.com/stretchr/testify/assert"
)

var component = "[API_HANDLER]"

func TestPing(t *testing.T) {
    expectedResult := api.SuccessResponse{Success: true, Response: "pong"}
    req, _ := http.NewRequest(http.MethodGet, "/api/debug/ping", nil)

    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, w.Code, 200)

    result := testutils.ParseAndValidateResponse(t, w)

    assert.Equal(t, result.Success, expectedResult.Success, "Response.Success shows false")
    assert.Equal(t, result.Response, expectedResult.Response)
}

func TestTruncate(t *testing.T) {
    expectedResult := true
    req, _ := http.NewRequest(http.MethodPost, "/api/debug/truncate", nil)
    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, w.Code, 200)

    result := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, result.Success, expectedResult, "Response.Success shows false: " + w.Body.String())
}

func TestHealthCheck(t *testing.T) {
    env := testutils.InitRouter(true)
    req, _ := http.NewRequest(http.MethodGet, "/api/debug/healthcheck", nil)
    res := env.SendTestingRequestSuccess(
        t,
        req,
        func() {},
        http.StatusOK,
        true,
    )

    result := res.Response.(map[string]interface{})
    assert.Equal(t, "successfull", result["databaseConnection"])
}
