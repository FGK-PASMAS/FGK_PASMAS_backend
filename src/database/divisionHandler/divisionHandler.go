package divisionhandler

import (
	"context"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	dberr "github.com/MetaEMK/FGK_PASMAS_backend/database/dbErr"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)


var log = logging.DivisionHandlerLogger

func GetDivisions() ([]model.DivisionStructSelect, error) {
    err := database.CheckDatabaseConnection()
    if err != nil {
        return []model.DivisionStructSelect{}, dberr.ErrNoConnection
    }

    query := `SELECT id, name, passenger_capacity FROM division`

    rows, err := database.PgConn.Query(context.Background(), query)
    defer rows.Close()

    if err != nil {
        errMessage := "Failed to got divisions from database"
        log.Info(errMessage)
        return []model.DivisionStructSelect{}, dberr.ErrQuery
    }

    divisions := []model.DivisionStructSelect{}
    for rows.Next() {
        var division model.DivisionStructSelect
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

func GetDivisionById(id int) (model.DivisionStructSelect, error) {
    err := database.CheckDatabaseConnection()
    if err != nil {
        return model.DivisionStructSelect{}, dberr.ErrNoConnection
    }

    query := `SELECT id, name, passenger_capacity FROM division WHERE id = $1`

    row := database.PgConn.QueryRow(context.Background(), query, id)

    var division model.DivisionStructSelect
    err = row.Scan(&division.Id, &division.Name, &division.PassengerCapacity)
    if err != nil {
        return model.DivisionStructSelect{}, err
    }

    return division, nil
}
