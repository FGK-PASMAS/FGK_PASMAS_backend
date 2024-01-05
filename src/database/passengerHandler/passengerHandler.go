package passengerhandler

import (
	"context"
	"fmt"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	passStream "github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes/passenger"
)


var log = logging.DbLogger

func GetPassengers() ([]PassengerStructSelect, error) {

    query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name FROM passenger p JOIN division d ON p.division_id = d.id`

    rows, err := database.PgConn.Query(context.Background(), query)
    if err != nil {
        log.Warn("Failed to got passengers from database: " + err.Error())
        return nil, err
    } else {
        log.Debug("Successfully got passengers from database")

        var passengers []PassengerStructSelect = make([]PassengerStructSelect, 0)

        for rows.Next() {
            var passenger PassengerStructSelect

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

func GetPassengerById(id int64) (PassengerStructSelect, error) {
    query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name FROM passenger p JOIN division d ON p.division_id = d.id WHERE p.id=$1`

    row := database.PgConn.QueryRow(context.Background(), query, id)

    var passenger PassengerStructSelect
    err := row.Scan(&passenger.Id, &passenger.LastName, &passenger.FirstName, &passenger.Weight, &passenger.Division.Id, &passenger.Division.Name)
    if err != nil {
        log.Warn(fmt.Sprintf("Failed to get passenger with id %d from database: %s", id, err.Error()))
        return PassengerStructSelect{}, err
    } else {
        log.Debug("Successfully got passenger from database")
        return passenger, nil
    }
}

func CreatePassenger(pass PassengerStructInsert) (PassengerStructSelect, error) {
    query := `INSERT INTO passenger (last_name, first_name, weight, division_id) VALUES ($1, $2, $3, $4) RETURNING id`

    res := database.PgConn.QueryRow(context.Background(), query, pass.LastName, pass.FirstName, pass.Weight, pass.DivisionId)

    var id int
    err := res.Scan(&id)
    if err != nil {
        log.Warn("Failed create passenger in database")
        return PassengerStructSelect{}, err
    }

    newPass, err := GetPassengerById(int64(id))
    if err != nil {
        log.Warn("Failed to get passenger from database")
        return PassengerStructSelect{}, err
    }

    log.Debug("Successfully created passenger in database")
    passStream.PublishPassengerEvent(realtime.CREATED, newPass)
    return newPass, nil
}

func UpdatePassenger(pass PassengerStructUpdate) (PassengerStructSelect, error) {
    query := `UPDATE passenger SET last_name=$2, first_name=$3, weight=$4, division_id=$5 WHERE id=$1 RETURNING id`

    res := database.PgConn.QueryRow(context.Background(), query, pass.Id, pass.LastName, pass.FirstName, pass.Weight, pass.DivisionId)

    var id int
    err := res.Scan(&id)
    if err != nil {
        log.Warn("Failed update passenger in database")
        return PassengerStructSelect{}, err
    }

    newPass, err := GetPassengerById(int64(id))
    if err != nil {
        log.Warn("Failed to get passenger from database")
        return PassengerStructSelect{}, err
    }

    log.Debug("Successfully updated passenger in database")
    passStream.PublishPassengerEvent(realtime.UPDATED, newPass)
    return newPass, nil
}

func DeletePassenger(id int) error {
    query := `DELETE FROM passenger WHERE id=$1`

    _, err := database.PgConn.Exec(context.Background(), query, id)
    if err != nil {
        log.Warn("Failed to delete passenger from database")
    }

    passStream.PublishPassengerEvent(realtime.DELETED, id)
    return err
}
