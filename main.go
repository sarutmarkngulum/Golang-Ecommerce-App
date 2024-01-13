package main

import (
	"os"

	"github.com/sarutmarkngulum/Golang-Ecommerce-App/modules/servers"

	"github.com/sarutmarkngulum/Golang-Ecommerce-App/config"
	"github.com/sarutmarkngulum/Golang-Ecommerce-App/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
