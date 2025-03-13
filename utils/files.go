package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

	return false
}


