package pasmasservice

import (
	"errors"

	dberr "github.com/MetaEMK/FGK_PASMAS_backend/database/dbErr"
)

//The following errors are returned by the database handler
var (
    //This error occurs when an unknown error occurs
    ErrUnknown = errors.New("UnknownError")

    // Error occurs when connection to database is lost
    ErrNoDbConnection = dberr.ErrNoConnection

    //Error occurs when a database query fails to execute
    ErrDbQuery = dberr.ErrQuery

    //Error occurs when no rows are returned - only for transactions into other Errors needed
    errDbNoRows = dberr.ErrNoRows

    // This error occurs when no object is found
    ErrObjectNotFound = errors.New("No object found")

    // This error occurs when the object already exists
    ErrObjectAlreadyExists = errors.New("Object already exists")

    ErrObjectCreatedFailed = errors.New("Object creation failed")

)


//The following errors occur when a dependency is missing
var(
    // This error occurs when the dependent division of the object are not found
    ErrObjectDependencyDivisionMissing = errors.New("Division does not exist")
) 
