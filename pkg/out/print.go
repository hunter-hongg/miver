package out

import (
	"fmt"
	"miver/pkg/color"
)

func Print(text string) {
	fmt.Println(text)
}

func Warning(text string) {
	fmt.Println(color.Yellow(fmt.Sprintf("Warning: %s", text)))
}

func Error(text string) {
	fmt.Println(color.Red(fmt.Sprintf("Error: %s", text)))
}

func Info(text string) {
	fmt.Println(color.Blue(fmt.Sprintf("Info: %s", text)))
}

func Step(text string) {
	fmt.Println(color.Cyan(fmt.Sprintf("Step: %s", text)))
}
