package errs

import (
	"miver/pkg/out"
	"os"
)

func DealError(err error) {
	if err != nil {
		out.Error(err.Error())
		os.Exit(2)
	}
}
