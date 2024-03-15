// ANSI escape codes
package ansi

const (
	osc = "\u001B]"
	bel = "\u0007"
	esc = "\u001B["
)

const (
	reset  = "0m"
	bold   = "1m"
	italic = "3m"
)

// make a hyperlink. https://gist.github.com/egmontkob/eb114294efbcd5adb1944c9f3cb5feda
func Link(text, uri string) string {
	return osc + "8;;" + uri + bel + text + osc + "8;;" + bel
}

func Bold(text string) string {
	return esc + bold + text + esc + reset
}

func Italic(text string) string {
	return esc + italic + text + esc + reset
}
