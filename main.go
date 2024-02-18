package main

import (
	"github.com/HakimovBulat/lazyCalculator/router"
	"github.com/HakimovBulat/lazyCalculator/utils"
	_ "github.com/lib/pq"
)

func main() {
	router := router.SetupRouter()
	if err := router.Run(":8080"); err != nil {
		utils.Logger.Error(err.Error())
	}
}
