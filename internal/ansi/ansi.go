// SPDX-FileCopyrightText: 2024 Bogdan Alekseevich Zazhigin <zaboal@tuta.io>
// SPDX-License-Identifier: 0BSD

// Package ansi provides escape, e.g. for stylish logs
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

func Link(text, uri string) string {
	return osc + "8;;" + uri + bel + text + osc + "8;;" + bel
}

func Bold(text string) string {
	return esc + bold + text + esc + reset
}

func Italic(text string) string {
	return esc + italic + text + esc + reset
}
