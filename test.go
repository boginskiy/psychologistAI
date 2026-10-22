package main

import (
	"fmt"
	"time"
)

const (
	HostEmail = "smtp.gmail.com"
	PortEmail = 587
	FromEmail = "gophkeeper@gmail.com"
	Password  = "upiplnvviujgnevc"

	Subject = "Подтверждение email"
	Retry   = 3
	TimeOut = 200 * time.Millisecond
)

func main() {

	for range 0 {
		fmt.Println("ddd")
	}

	// _ = db
}
