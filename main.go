package main

import (
	"date-app/config"
	"date-app/handler/route"
	"date-app/utils/jwt"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	err := config.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}
	err = config.ConnectMySQL()
	if err != nil {
		log.Fatal(err)
	}
	jwt.Tokens = map[string]bool{}
	router := route.NewRoute()
	err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router)
	if err != nil {
		log.Fatal(err)
	}
}
