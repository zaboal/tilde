package main

const (
	osc = "\u001B]"
	bel = "\u0007"
)

// создание ссылки с помощью osc8: gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func link(uri, text string) string {
	return osc + "8;;" + uri + bel + text + osc + "8;;" + bel
}
