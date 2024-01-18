package pasmasservice

import (
	dberr "github.com/MetaEMK/FGK_PASMAS_backend/database/dbErr"
	divisionhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/divisionHandler"
	passengerhandler "github.com/MetaEMK/FGK_PASMAS_backend/database/passengerHandler"
	"github.com/MetaEMK/FGK_PASMAS_backend/model"
)

func GetPassengers() ([]model.PassengerStructSelect, error) {
    return passengerhandler.GetPassengers()
}

func CreatePassenger(pass model.PassengerStructInsert) (model.PassengerStructSelect, error) {
    // Check for dependencies
    // Check for division
    _, err := divisionhandler.GetDivisionById(pass.DivisionId)
    if err != nil {
        if err == dberr.ErrNoRows {
            return model.PassengerStructSelect{}, ErrObjectDependencyDivisionMissing
        } else {
            return model.PassengerStructSelect{}, err
        }
    }


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

    // Check for division
    _, err = divisionhandler.GetDivisionById(pass.DivisionId)
    if err != nil {
        if err == dberr.ErrNoRows {
            return model.PassengerStructSelect{}, ErrObjectDependencyDivisionMissing
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

    return nil
}