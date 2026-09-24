package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/oschwald/geoip2-golang"
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

	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		fmt.Println("не удалось открыть базу GeoIP")
		// return false, errors.New("не удалось открыть базу GeoIP")
	}

	ip := net.ParseIP("8.8.8.8")

	record, err := db.Country(ip)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(record.Country.IsoCode)

	// _ = db
}
