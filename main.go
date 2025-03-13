package main

import (
	"fmt"

	"decision-tree/utils"
)

func main() {
	utils.ParseFlag()
	if *utils.CommandPtr != "train" && *utils.CommandPtr != "predict" {
		fmt.Println("Please provide a valid command")
		fmt.Println("Ex: -c train or -c predict")
		return
	}
	if *utils.InputPtr == "" {
		fmt.Println("Please provide an input file")
		fmt.Println("Ex: -i <filepath.csv>")
		return
	}
	if *utils.ColumnPtr == "" && *utils.CommandPtr == "train" {
		fmt.Println("Please provide a column to train")
		fmt.Println("Ex: -t <column_name>")
		return
	}
	if *utils.ModelFilePtr == "" && *utils.CommandPtr == "predict" {
		fmt.Println("Please provide a trained decision tree to predict")
		fmt.Println("Ex: -m <filepath.dt>")
		return
	}
}
