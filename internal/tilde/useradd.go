package tilde

import (
	"log"
	"os"
	"os/exec"
	"time"

	sha512 "github.com/GehirnInc/crypt/sha512_crypt"
	sethvargo "github.com/sethvargo/go-password/password"

	ansi "github.com/zaboal/tilde/internal/ansi"
)

var logger = log.New(os.Stdout, ansi.Bold("tilde "), log.Ldate+log.Ltime+log.Lmsgprefix)

// generate and hash a password as for /etc/shadow
func Password() (password string, shadow string) {
	password, _ = sethvargo.Generate(8, 4, 0, true, true)
	shadow, _ = sha512.New().Generate([]byte(password), nil)

	return password, shadow
}

// subcribe a user onto tilde
func Subscribe(username string) (password string, error error) {
	password, shadow := Password()
	expire := time.Now().AddDate(0, 3, 0).Format("2006-01-02")

	command := exec.Command("useradd",
		"--groups", "subcribers",
		"--create-home",
		"--inactive", "7",
		"--expiredate", expire,
		"--password", shadow,
		username)

	logger.Printf("adds login \"%s\" with shadow \"%s\"", username, shadow)
	return password, command.Run()
}
