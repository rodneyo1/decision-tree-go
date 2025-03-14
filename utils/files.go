package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"decision-tree/models"
)

func FileExtValidation() error {
	inputExt := filepath.Ext(*InputPtr)
	if *CommandPtr == "train" && inputExt != ".csv" {
		return errors.New("input file must be a CSV for training")
	}
	if *CommandPtr == "predict" && inputExt != ".csv" {
		return errors.New("input file must be a CSV for prediction")
	}
	if *CommandPtr == "train" && filepath.Ext(*OutputPtr) != ".dt" {
		return errors.New("output file must have .dt extension for model")
	}
	if *CommandPtr == "predict" && filepath.Ext(*OutputPtr) != ".csv" {
		return errors.New("output file must have .csv extension for predictions")
	}
	if *CommandPtr == "predict" && filepath.Ext(*ModelFilePtr) != ".dt" {
		return errors.New("model file must have .dt extension")
	}
	return nil
}

func ContainsClassAttribute() bool {
	csvFile, err := os.Open(*InputPtr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	columns, err := csvReader.Read()
	models.Columns = columns

	if err != nil {
		fmt.Println("failed to read file content")
		return false
	}
	if !slices.Contains(columns, *ColumnPtr) {
		return false
	}
	var batch []map[string]interface{}

	for {
		row, err := csvReader.Read()

		if err == io.EOF {
			break
		}
		record := make(map[string]interface{})
		for i, val := range row {
			record[columns[i]] = parseValue(val)
		}
		batch = append(batch, record)

		if len(batch) == 100 {
			// train
			batch = nil
		}

	}
	fmt.Println(batch[0])

	if len(batch) > 0 {
		//train

	}

	return false
}
