package main

import (
	"fmt"
	"os"

	"dt/algorithm"
	"dt/utils"
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
	if *utils.OutputPtr == "" {
		fmt.Println("Please provide an output file")
		fmt.Println("Ex: -o <filepath.dt> for training or -o <filepath.csv> for prediction")
		return
	}
	err := utils.FileExtValidation()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if *utils.CommandPtr == "train" {
		err = runTraining()
	} else if *utils.CommandPtr == "predict" {
		err = runPrediction()
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// runTraining handles the training workflow
func runTraining() error {
	fmt.Println("Starting training process...")
	if err := utils.LoadTrainingData(); err != nil {
		return fmt.Errorf("failed to load training data: %w", err)
	}
	// Build the decision tree
	tree, err := algorithm.BuildTree(*utils.ColumnPtr)
	if err != nil {
		return fmt.Errorf("failed to build decision tree: %w", err)
	}

	// Save the model
	if err := utils.SaveModel(tree); err != nil {
		return fmt.Errorf("failed to save model: %w", err)
	}

	fmt.Println("Training completed successfully!")
	return nil
}

// runPrediction handles the prediction workflow
func runPrediction() error {
	fmt.Println("Starting prediction process...")

	// Load the model
	modelData, err := utils.LoadModels()
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}
	if err := utils.LoadPredictionData(); err != nil {
		return fmt.Errorf("failed to load prediction data: %w", err)
	}

	// Make predictions
	predictions := algorithm.Predict(modelData.Tree)
	// Save predictions
	if err := utils.SavePredictions(predictions); err != nil {
		return fmt.Errorf("failed to save predictions: %w", err)
	}

	fmt.Println("Prediction completed successfully!")
	return nil
}
