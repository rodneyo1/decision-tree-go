package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"decision-tree/models"
)

func FileExtValidation() error {
	if !(filepath.Ext(*InputPtr) != ".csv" || filepath.Ext(*InputPtr) == ".dt") {
		fmt.Println("invalid file format. Either use .csv for input and .dt for output")
		return errors.New("invalid file format")
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
			for i, val := range row{
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

func parseValue(value string) interface{} {
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	if value == "true" || value == "false" {
		return value == "true"
	}

	if dateVal, err := time.Parse("2006-01-02", value); err == nil {
		return dateVal
	}

	return value
}
