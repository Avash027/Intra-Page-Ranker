package logger

import (
	"fmt"
)

var ColorRed = "\033[31m"
var ColorGreen = "\033[32m"
var ColorYellow = "\033[33m"
var ColorBlue = "\033[34m"
var ColorPurple = "\033[35m"
var ColorCyan = "\033[36m"
var ColorWhite = "\033[37m"

func Red(msg string) {
	fmt.Println(string(ColorRed), msg, string(ColorWhite))
}

func Green(msg string) {
	fmt.Println(string(ColorGreen), msg, string(ColorWhite))
}

func Yellow(msg string) {
	fmt.Println(string(ColorYellow), msg, string(ColorWhite))
}

func Blue(msg string) {
	fmt.Println(string(ColorBlue), msg, string(ColorWhite))
}

func Purple(msg string) {
	fmt.Println(string(ColorPurple), msg, string(ColorWhite))
}

func Cyan(msg string) {
	fmt.Println(string(ColorCyan), msg, string(ColorWhite))
}

func InPlaceCyan(msg string) {
	fmt.Print("\r", string(ColorCyan), msg, string(ColorWhite))
}
