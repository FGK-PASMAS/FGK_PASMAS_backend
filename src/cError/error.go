package cerror

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
    ErrPassengerWeightIsZero = errors.New("Passenger weight is zero")
    ErrObjectDependencyMissing = errors.New("Object dependency missing")
    ErrObjectNotFound = gorm.ErrRecordNotFound
    ErrNoRealtimeHandlerFound = errors.New("No realtime handler found")
    ErrIncludeNotSupported = errors.New("Include not supported")
    ErrFilterNotSupported = errors.New("Filter not supported")

    ErrRealtimeEventCouldNotBeCreated = errors.New("Realtime event could not be created")
    ErrPassengerActionNotValid = errors.New("Passenger action not valid")
)

// Internal Errors only - this should never occur
var (
    ErrDatabaseHandlerDestroy = errors.New("Struct DatabaseHandler was never closed correctly")
)

// Authentication and Authorisation
var (
    ErrPasswordTooLong = bcrypt.ErrPasswordTooLong
    ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword
    ErrInvalidCredentials = errors.New("Invalid credentials")
    ErrEmptyCredentials = errors.New("Empty token, password or username")
)
