package utils

import "flag"

var (
	CommandPtr   = flag.String("c", "", "Specify the command")
	InputPtr     = flag.String("i", "", "input csv file")
	ColumnPtr    = flag.String("t", "", "name of the target column")
	OutputPtr    = flag.String("o", "", "path to save trained dataset tree model")
	ModelFilePtr = flag.String("m", "", "path to trained dataset for predictions")
)

func ParseFlag() {
	flag.Parse()
}
