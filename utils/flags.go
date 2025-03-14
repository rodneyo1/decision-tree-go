package utils

import "flag"

var (
	CommandPtr   = flag.String("c", "", "Specify the command")
	InputPtr     = flag.String("i", "", "input csv file")
	ColumnPtr    = flag.String("t", "", "name of the target column")
	OutputPtr    = flag.String("o", "", "path to save trained dataset tree model")
	ModelFilePtr = flag.String("m", "", "path to trained dataset for predictions")
	BinningEnabledPtr = flag.Bool("bin", false, "Enable binning for numerical features")
    BinningMethodPtr  = flag.String("binmethod", "equal_width", "Binning method (equal_width or equal_frequency)")
    NumBinsPtr        = flag.Int("numbins", 10, "Number of bins to use")
)

func ParseFlag() {
	flag.Parse()
}
