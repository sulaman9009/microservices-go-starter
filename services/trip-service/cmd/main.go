package main

import (
	"log"
	"time"
)

func main() {
	for {
		time.Sleep(10 * time.Second)
		log.Println("Trip service is running...")
	}
}
