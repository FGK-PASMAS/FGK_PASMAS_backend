package main

import (
	"github.com/MetaEMK/FGK_PASMAS_backend/router"
)

func main() {
    server := router.InitRouter() 
    server.Run(":8080")
}
