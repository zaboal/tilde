package main

import (
	"log"
	"os/exec"
	"time"

	sha512 "github.com/GehirnInc/crypt/sha512_crypt"
	sethvargo "github.com/sethvargo/go-password/password"
)

// регистрация пользователя в тильде забоала
func register(username string) (password string, error error) {
	password, _ = sethvargo.Generate(8, 4, 0, true, true)
	shadow, _ := sha512.New().Generate([]byte(password), nil)

	// добавить три месяца к текущей дате и вывести результат в формате iso
	expire := time.Now().AddDate(0, 3, 0).Format("2006-01-02")

	command := exec.Command("useradd",
		"--groups", "subcribers",
		"--create-home",
		"--inactive", "7",
		"--expiredate", expire,
		"--password", shadow,
		username)

	log.Printf("executing \"%s\"", command.String())
	return password, command.Run()
}
