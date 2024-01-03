package passengerhandler

import (
	"context"
	"fmt"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
)


var log = logging.DbLogger

func GetPassengers() ([]SelectPassenger, error) {

    query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name FROM passenger p JOIN division d ON p.division_id = d.id`

    rows, err := database.PgConn.Query(context.Background(), query)
    if err != nil {
        log.Warn("Failed to got passengers from database: " + err.Error())
        return nil, err
    } else {
        log.Debug("Successfully got passengers from database")

        var passengers []SelectPassenger = make([]SelectPassenger, 0)

        for rows.Next() {
            var passenger SelectPassenger

            err = rows.Scan(&passenger.Id, &passenger.LastName, &passenger.FirstName, &passenger.Weight, &passenger.Division.Id, &passenger.Division.Name)
            if(err != nil) {
                log.Warn("Could not parse passenger row from database to passenger type - skipping entry")
                log.Debug(err.Error())
            } else {
                passengers = append(passengers, passenger)
            }
        }
        return passengers, nil
    }
}

func GetPassengerById(id int64) (SelectPassenger, error) {
    query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name FROM passenger p JOIN division d ON p.division_id = d.id WHERE p.id=$1`

    row := database.PgConn.QueryRow(context.Background(), query, id)

    var passenger SelectPassenger
    err := row.Scan(&passenger.Id, &passenger.LastName, &passenger.FirstName, &passenger.Weight, &passenger.Division.Id, &passenger.Division.Name)
    if err != nil {
        log.Warn(fmt.Sprintf("Failed to get passenger with id %d from database: %s", id, err.Error()))
        return SelectPassenger{}, err
    } else {
        log.Debug("Successfully got passenger from database")
        return passenger, nil
    }
}

func CreatePassenger(pass InsertPassenger) (DatabasePassenger, error) {
    query := `INSERT INTO passenger (last_name, first_name, weight, division_id) VALUES ($1, $2, $3, $4) RETURNING id, last_name, first_name, weight, division_id`

    res := database.PgConn.QueryRow(context.Background(), query, pass.LastName, pass.FirstName, pass.Weight, pass.DivisionId)

    var newPass DatabasePassenger
    err := res.Scan(&newPass.Id, &newPass.LastName, &newPass.FirstName, &newPass.Weight, &newPass.DivisionId)
    if err != nil {
        log.Warn("Failed create passenger in database")
        return DatabasePassenger{}, err
    }

    log.Debug("Successfully created passenger in database")
    return newPass, nil
}

func UpdatePassengerById(id int64, pass InsertPassenger) {
    //query := `UPDATE passenger SET last_name=$2, first_name=$3, weight=$4, division_id$5 WHERE id=$1`
}
