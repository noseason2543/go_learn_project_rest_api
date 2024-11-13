package main

import (
	"fmt"
	"go_learn_project_rest_api/config"
	"go_learn_project_rest_api/pkgs/databases"
	"os"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	fmt.Println("hello world")
	cfg := config.LoadConfig(envPath())
	fmt.Println(cfg.App().Url())

	db := databases.DbConnection(cfg.Db())
	defer db.Close()
	fmt.Println("db connect : ", db)
}
