package passengerhandler

import (
	"context"

	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	dberr "github.com/MetaEMK/FGK_PASMAS_backend/database/dbErr"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	passStream "github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes/passenger"
)

var log = logging.PassHandlerLogger


func GetPassengers() ([]model.PassengerStructSelect, error) {
    connErr := database.CheckDatabaseConnection()
    if connErr != nil {
        return []model.PassengerStructSelect{}, connErr
    }

	query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name, d.passenger_capacity FROM passenger p JOIN division d ON p.division_id = d.id;`

	rows, err := database.PgConn.Query(context.Background(), query)
	if err != nil {
        errMessage := "Failed to got passengers from database"
		log.Info(errMessage + " - " + err.Error())

        return nil, dberr.ErrQuery
	} else {
		log.Debug("Successfully got passengers from database")

		var passengers []model.PassengerStructSelect = make([]model.PassengerStructSelect, 0)

		for rows.Next() {
			var passenger model.PassengerStructSelect

			err = rows.Scan(
                &passenger.Id,
                &passenger.LastName,
                &passenger.FirstName,
                &passenger.Weight,
                &passenger.Division.Id,
                &passenger.Division.Name,
                &passenger.Division.PassengerCapacity,
            )
			if err != nil {
				log.Info("Could not parse passenger row from database to passenger type - skipping entry")
				log.Debug(err.Error())
			} else {
				passengers = append(passengers, passenger)
			}
		}
        return passengers, nil
	}
}

func GetPassengerById(id int64) (model.PassengerStructSelect, error) {
    connErr := database.CheckDatabaseConnection()
    if connErr != nil {
        return model.PassengerStructSelect{}, dberr.ErrNoConnection
    }

	query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name, d.passenger_capacity FROM passenger p JOIN division d ON p.division_id = d.id WHERE p.id=$1`

	row := database.PgConn.QueryRow(context.Background(), query, id)

	var passenger model.PassengerStructSelect
	err := row.Scan(
        &passenger.Id,
        &passenger.LastName,
        &passenger.FirstName,
        &passenger.Weight,
        &passenger.Division.Id,
        &passenger.Division.Name,
        &passenger.Division.PassengerCapacity,
    )

	if err != nil {
        return model.PassengerStructSelect{}, err
	} else {
		log.Debug("Successfully got passenger from database")
		return passenger, nil
	}
}

func CreatePassenger(pass model.PassengerStructInsert) (int64, error) {
    connErr := database.CheckDatabaseConnection()
    if connErr != nil {
        return -1, dberr.ErrNoConnection
    }

	query := `INSERT INTO passenger (last_name, first_name, weight, division_id) VALUES ($1, $2, $3, $4) RETURNING id`

	res := database.PgConn.QueryRow(context.Background(), query, pass.LastName, pass.FirstName, pass.Weight, pass.DivisionId)

	var id int64
	err := res.Scan(&id)

	if err != nil {
		return -1, dberr.ErrQuery
	}

    log.Debug("Successfully created passenger in database")
    return id, nil
}

// UpdatePassenger updates a passenger in the database
// Returns the id of the updated passenger
// Returns an error if the passenger could not be updated
func UpdatePassenger(id int64, pass model.PassengerStructUpdate) (int64, error) {
    connErr := database.CheckDatabaseConnection()
    if connErr != nil {
        return -1, dberr.ErrNoConnection
    }
	query := `UPDATE passenger SET last_name=$2, first_name=$3, weight=$4, division_id=$5 WHERE id=$1 RETURNING id`

	res := database.PgConn.QueryRow(context.Background(), query, id, pass.LastName, pass.FirstName, pass.Weight, pass.DivisionId)

	var newId int64
	err := res.Scan(&newId)
	if err != nil {
        log.Debug("Failed to update passenger: + " + err.Error())
		return -1, dberr.ErrQuery
	}

    log.Debug("Successfully updated passenger in database")
    return newId, nil
}

func DeletePassenger(id int64) error {
    connErr := database.CheckDatabaseConnection()
    if connErr != nil {
        return dberr.ErrNoConnection
    }

	query := `DELETE FROM passenger WHERE id=$1`

	_, err := database.PgConn.Exec(context.Background(), query, id)
	if err != nil {
        return dberr.ErrQuery
	} else {
        log.Debug("Successfully deleted passenger from database")
        passStream.PublishPassengerEvent(realtime.DELETED, id)
        return nil
    }
}
