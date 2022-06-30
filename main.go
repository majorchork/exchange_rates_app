package main

import (
	"github.com/majorchork/rates_app/server"
	"log"
)

func main() {
	// what is handler? should be server
	log.Println("server run")
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("server started")

}
