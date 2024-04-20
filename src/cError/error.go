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
    ErrForbidden = errors.New("You are not allowed to access this resource")
    ErrInvalidRole = errors.New("Invalid role")
    ErrUserAlreadyExists = errors.New("User already exists")
)

var (
    ErrTooMuchFuel = errors.New("Too much fuel")
    ErrTooLessFuel = errors.New("Too less fuel")
    ErrSlotIsNotFree = errors.New("Slot is not free")
    ErrFlightStatusDoesNotFitProcess = errors.New("Flight status does not fit process")
    ErrDepartureTimeIsZero = errors.New("Departure time is zero")
    ErrInvalidArrivalTime = errors.New("Invalid arrival time")
    ErrNoPilotAvailable = errors.New("No valid pilot available")
    ErrNoStartFuelFound = errors.New("No start fuel found")
    ErrMaxSeatPayload = errors.New("maxSeatPayload was exceeded")
    ErrTooManyPassenger = errors.New("too many passengers for this plane")
    ErrTooLessPassenger = errors.New("A flight needs to have at least one passenger")
    ErrOverloaded = errors.New("MTOW is exceeded")
    ErrFlightNoCouldNotBeGenerated = errors.New("Could not generate flightNo")
)
