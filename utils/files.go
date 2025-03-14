package utils

import (
	"errors"
	"path/filepath"
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
