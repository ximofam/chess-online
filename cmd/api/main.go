package main

import (
	"fmt"
	"log"

	"github.com/ximofam/chess-online/app"
	"github.com/ximofam/chess-online/config"
	"github.com/ximofam/chess-online/database"
	"github.com/ximofam/chess-online/services/auth"
	"github.com/ximofam/chess-online/services/auth/token"
)

func main() {
	envs := config.Envs
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC",
		"root", envs.MySQL.DBPassword, envs.MySQL.DBHost, envs.MySQL.DBPort, envs.MySQL.DBName,
	)

	db, err := database.New(database.Config{
		Driver: "mysql",
		DSN:    dsn,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	app := app.New(db)
	token.Init(envs.JWT)
	auth.Init(db)

	app.Run()
}
