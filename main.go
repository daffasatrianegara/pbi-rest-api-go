package main 

import (
	"task-5-pbi-btpns-daffa_satria/config"
	"task-5-pbi-btpns-daffa_satria/router"
)
func main() {
	config.Connect()
	config.Migrate()
	r := router.SetupRouter()
	r.Run(":3000")


}