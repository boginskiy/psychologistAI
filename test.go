package main

import (
	"fmt"

	"github.com/boginskiy/psychologistAI/pkg/hashpass"
)

func main() {
	pass := "1234"
	hPass, _ := hashpass.CreateHashPass(pass)

	err := hashpass.CheckPassword(hPass, pass)
	if err == nil {
		fmt.Println("password is correct")
	}
}
