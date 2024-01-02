package passengerhandler

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/database"
	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
)


var log = logging.DbLogger

func GetPassengers() ([]Passenger, error) {

    checkDb, dbErr := database.CheckDatabaseConnection()
    if checkDb == false {
        return nil, dbErr
    }

    query := `SELECT p.id, p.last_name, p.first_name, p.weight, d.id, d.name FROM passenger p JOIN division d ON p.division_id = d.id`

    rows, err := database.Db.Query(query)
    if err != nil {
        log.Warn("Failed to got passengers from database")
        return nil, err
    } else {
        log.Debug("Successfully got passengers from database")

        var passengers []Passenger = make([]Passenger, 0)

        for rows.Next() {
            var passenger Passenger

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

