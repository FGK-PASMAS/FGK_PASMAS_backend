package noGen

import (
	"fmt"
	"strconv"
	"strings"

	cerror "github.com/MetaEMK/FGK_PASMAS_backend/cError"
	databasehandler "github.com/MetaEMK/FGK_PASMAS_backend/databaseHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GenerateFlightNo(plane model.Plane) (string, error) {
    var flightNo string

    var prevFlight model.Flight
    err := databasehandler.Db.Unscoped().Model(model.Flight{}).Where("plane_id = ?", plane.ID).Where("flight_no IS NOT NULL").Order("flight_no DESC").First(&prevFlight).Error

    if err != nil {
        if err == cerror.ErrObjectNotFound {
            flightNo = generateFlightNumberPattern(plane, 1)
            println(flightNo)
            return flightNo, nil
        } else {
            println(err.Error())
            return "", cerror.ErrFlightNoCouldNotBeGenerated
        }
    }

    n, err := getFlightNumber(*prevFlight.FlightNo)
    n++

    flightNo = generateFlightNumberPattern(plane, n)

    return flightNo, err
}


func generateFlightNumberPattern(plane model.Plane, number uint) string {
    reg := strings.Split(plane.Registration, "-")
    prefix := reg[len(reg)-1]

    return fmt.Sprintf("%s-%03d", prefix, number)
}

func getFlightNumber(flightNo string) (uint, error) {
    strs := strings.Split(flightNo, "-")
    no := strs[len(strs)-1]
    n, err := strconv.ParseUint(no, 10, 64)

    return uint(n), err
}
