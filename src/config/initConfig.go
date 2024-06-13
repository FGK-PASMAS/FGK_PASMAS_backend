package config

import (
	"os"
	"strconv"

	"github.com/MetaEMK/FGK_PASMAS_backend/logging"
)

var log = logging.NewLogger("CONFIG", GetGlobalLogLevel())

func GetGlobalLogLevel() logging.LogLevel {
    logLevelStr := os.Getenv("GLOBALLOGLEVEL")
    logLevel, err := strconv.Atoi(logLevelStr)
    if err != nil {
        return logging.INFO
    }

    return logging.LogLevel(logLevel)
}
