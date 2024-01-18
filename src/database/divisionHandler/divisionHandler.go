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
    err := database.CheckDatabaseConnection()
    if err != nil {
        return []DivisionStructSelect{}, intError{Type: internalerror.ErrorDatabaseConnectionError, Message: "Failed to connect to database", Body: err}
    }

    query := `SELECT id, name, passenger_capacity FROM division`

    rows, err := database.PgConn.Query(context.Background(), query)
    defer rows.Close()

    if err != nil {
        errMessage := "Failed to got divisions from database"
        log.Info(errMessage)
        return []DivisionStructSelect{}, intError{Type: internalerror.ErrorDatabaseQueryError, Message: errMessage, Body: err}
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

func GetDivisionById(id int) (DivisionStructSelect, error) {
    err := database.CheckDatabaseConnection()
    if err != nil {
        return DivisionStructSelect{}, intError{Type: internalerror.ErrorDatabaseConnectionError, Message: "Failed to connect to database", Body: err}
    }

    query := `SELECT id, name, passenger_capacity FROM division WHERE id = $1`

    row := database.PgConn.QueryRow(context.Background(), query, id)

    var division DivisionStructSelect
    err = row.Scan(&division.Id, &division.Name, &division.PassengerCapacity)
    if err != nil {
        errMessage := "Failed to parse one division from database - skipping"
        log.Info(errMessage)
        return DivisionStructSelect{}, intError{Type: internalerror.ErrorDatabaseQueryError, Message: errMessage, Body: err}
    }

    return division, nil
}
