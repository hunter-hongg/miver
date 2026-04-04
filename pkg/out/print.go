package out

import (
	"fmt"
	"miver/pkg/color"
)

func Print(text string) {
	fmt.Println(text)
}

func Warning(text string) {
	fmt.Println(color.Yellow("Warning"), text)
}

func Error(text string) {
	fmt.Println(color.Red("Error"), text)
}

func Info(text string) {
	fmt.Println(color.Blue("Info"), text)
}

func Step(text string) {
	fmt.Println(color.Cyan("Step"), text)
}
