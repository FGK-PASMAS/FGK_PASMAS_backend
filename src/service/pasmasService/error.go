package pasmasservice

import (
	"errors"

	"gorm.io/gorm"
)

//The following errors are returned by the database handler
var (
    //This error occurs when an unknown error occurs
    ErrUnknown = errors.New("UnknownError")

    ErrDataBaseErr = errors.New("An error with the database occured")

    // This error occurs when no object is found
    ErrObjectNotFound = gorm.ErrRecordNotFound

    // This error occurs when the object already exists
    ErrObjectAlreadyExists = errors.New("Object already exists")

    ErrObjectCreatedFailed = errors.New("Object creation failed")

)


//The following errors occur when a dependency is missing
var(
    // This error occurs when the dependent division of the object are not found
    ErrObjectDependencyMissing = errors.New("A dependent ressource does not exist")
    ErrIncludeNotSupported = errors.New("Include not supported")
) 


var (
    ErrTooLessFuel = errors.New("Not enough fuel for this flight")
    ErrTooMuchFuel = errors.New("More fuel than the plane can hold")
)
