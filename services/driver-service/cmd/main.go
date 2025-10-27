package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("driver service is running...")
		time.Sleep(60 * time.Second)
	}
}
