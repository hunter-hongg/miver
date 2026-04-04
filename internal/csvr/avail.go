package csvr

import (
	"encoding/csv"
	"fmt"
	"miver/pkg/errs"
	"os"
)

func ReadAvailable(miva_repo string) [][]string {
	csv_file := fmt.Sprintf("%s/available.csv", miva_repo)
	csv_fileh, err := os.Open(csv_file)
	errs.DealError(err)
	defer csv_fileh.Close()
	csv_reader := csv.NewReader(csv_fileh)
	lists, err := csv_reader.ReadAll()
	errs.DealError(err)
	return lists
}
