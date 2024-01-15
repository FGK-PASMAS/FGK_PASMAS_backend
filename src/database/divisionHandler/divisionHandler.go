package divisionhandler

import (
	"context"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	internalerror "github.com/MetaEMK/FGK_PASMAS_backend/internalError"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
)


var log = logging.DivisionHandlerLogger
type intError = internalerror.InternalError

func GetDivision() ([]DivisionStructSelect, error) {
    query := `SELECT id, name, passenger_capacity FROM division`

    rows, err := database.PgConn.Query(context.Background(), query)
    //defer rows.Close()

    if err != nil {
        errMessage := "Failed to got divisions from database"
        log.Info(errMessage)
        return []DivisionStructSelect{}, intError{Type: internalerror.DatabaseQueryError, Message: errMessage, Body: err}
    }

    divisions := []DivisionStructSelect{}
    for rows.Next() {
        var division DivisionStructSelect
        err = rows.Scan(&division.Id, &division.Name, &division.PassengerCapacity)
        if err != nil {
            errMessage := "Failed to parse one division from database - skipping"
            log.Info(errMessage)
        } else {
            divisions = append(divisions, division)
        }
    }

    return divisions, nil
}
