package utils

type ConsoleColor string

const (
	Red     ConsoleColor = "red"
	Green   ConsoleColor = "green"
	Yellow  ConsoleColor = "yellow"
	Blue    ConsoleColor = "blue"
	Magenta ConsoleColor = "magenta"
	Cyan    ConsoleColor = "cyan"
	White   ConsoleColor = "white"
	Default ConsoleColor = "default"
)

var colors = map[ConsoleColor]string{
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	White:   "\033[37m",
	Default: "\033[0m",
}

func GetColor(color ConsoleColor) string {
	return colors[color]
}

func Paint(text string, color ConsoleColor) string {
	return colors[color] + text + colors[Default]
}
