package pasmasservice

import (
	dberr "github.com/MetaEMK/FGK_PASMAS_backend/database/dbErr"
	passengerhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime"
	"github.com/MetaEMK/FGK_PASMAS_backend/router/realtime/routes/passenger"
)

func GetPassengers() ([]model.PassengerStructSelect, error) {
    return passengerhandler.GetPassengers()
}

func CreatePassenger(pass model.PassengerStructInsert) (model.PassengerStructSelect, error) {
    // Check for dependencies

    // Create passenger
    id, err := passengerhandler.CreatePassenger(pass)
    if err != nil {
        return model.PassengerStructSelect{}, err
    }

    //Get Passenger
    newPass, err := passengerhandler.GetPassengerById(id)
    if err != nil {
        return model.PassengerStructSelect{}, ErrObjectCreatedFailed
    }

    passenger.PublishPassengerEvent(realtime.CREATED, newPass)
    return newPass, nil
}

func UpdatePassenger(id int64, pass model.PassengerStructUpdate) (model.PassengerStructSelect, error) {
    // Check for dependencies
    // Check for old passenger
    _, err := passengerhandler.GetPassengerById(id)
    if err != nil {
        if(err == dberr.ErrNoRows) {
            return model.PassengerStructSelect{}, ErrObjectNotFound
        } else {
            return model.PassengerStructSelect{}, err
        }
    }

    // Update passenger
    newId, err := passengerhandler.UpdatePassenger(id, pass)
    if err != nil {
        return model.PassengerStructSelect{}, err
    }

    //Get GetPassenger
    newPass, err := passengerhandler.GetPassengerById(newId)
    if err != nil {
        return model.PassengerStructSelect{}, ErrObjectCreatedFailed
    }

    passenger.PublishPassengerEvent(realtime.UPDATED, newPass)
    return newPass, nil
}

func DeletePassenger(id int64) error {
    // Check for dependencies
    // Check for passenger to delete
    _, err := passengerhandler.GetPassengerById(id)
    if err != nil {
        if(err == dberr.ErrNoRows) {
            return ErrObjectNotFound
        } else {
            return err
        }
    }


    // Delete passenger
    err = passengerhandler.DeletePassenger(id)
    if err != nil {
        return err
    }

    passenger.PublishPassengerEvent(realtime.DELETED, id)
    return nil
}
