package main

import (
	"ohyandex/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		panic("No .env.local file")
	}
	http := server.NewHttp()
	http.Start()
}
