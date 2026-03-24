package color

func Colorize(text string, color string) string {
	return color + text + "\033[0m"
}

func Red(text string) string {
	return Colorize(text, "\033[31m")
}

func Green(text string) string {
	return Colorize(text, "\033[32m")
}

func Yellow(text string) string {
	return Colorize(text, "\033[33m")
}

func Blue(text string) string {
	return Colorize(text, "\033[34m")
}

func Magenta(text string) string {
	return Colorize(text, "\033[35m")
}

func Cyan(text string) string {
	return Colorize(text, "\033[36m")
}

func White(text string) string {
	return Colorize(text, "\033[37m")
}
