package logging

var DbLogger = NewLogger("DATABASE", DEBUG)
var PassHandlerLogger = NewLogger("PASSENGER_HANDLER", DEBUG)
var DivisionHandlerLogger = NewLogger("DIVISION_HANDLER", DEBUG)
var ApiLogger = NewLogger("API", DEBUG)
var DbDebugLogger = NewLogger("DB_DEBUG", DEBUG)

