package division_test

import (
	"encoding/json"
	"net/http"
	"testing"

	testutils "github.com/MetaEMK/FGK_PASMAS_backend/testUtils"
	"github.com/stretchr/testify/assert"
)


func TestGetDivision(t *testing.T) {
    req, _ := http.NewRequest(http.MethodGet, "/api/division/", nil)
    w := testutils.SendTestingRequest(t, req)

    assert.Equal(t, http.StatusOK, w.Code)

    res := testutils.ParseAndValidateResponse(t, w)
    assert.Equal(t, true, res.Success)

    var divisions []testutils.DivisionModel
    jsonBytes, _ := json.Marshal(res.Response)

    err := json.Unmarshal(jsonBytes, &divisions)
    assert.Nil(t, err)

    assert.Equal(t, 3, len(divisions))

    for _, division := range divisions {
        testutils.ValidateDivision(t, division)
    }
}
