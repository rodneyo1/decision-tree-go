# Fast & Scalable Decision Tree (C4.5) in Go

## Overview

This project is an implementation of the **C4.5 Decision Tree** classifier in Go. It is designed to be **high-performance**, **scalable**, and capable of handling **large datasets** with minimal memory overhead. The implementation supports **categorical and numerical attributes**, **missing value handling**, and **parallelization** for improved efficiency.

## Features

- **Train & Predict**: CLI commands for training a decision tree and making predictions.
- **Scalability**: Efficient handling of large datasets with optimized memory usage.
- **Modular Design**: Well-structured code for easy maintainability and extension.
- **Error Handling**: Clear error messages and validation for input/output files.
- **JSON Model Serialization**: Stores trained models in a JSON format.

## Installation

Ensure you have **Go 1.18+** installed. Then, clone the repository and build the executable:

```sh
git clone https://learn.zone01kisumu.ke/git/rodnochieng/decision-tree.git
cd decision-tree-go
go build -o dt
```

## Usage

### 1. Training a Decision Tree

```sh
./dt -c train -i <input_data_file.csv> -t <target_column> -o <output_tree.dt>
```

**Arguments:**
- `-c train` → Specifies the training command.
- `-i <input_data_file.csv>` → Path to the training dataset (CSV).
- `-t <target_column>` → Column name containing target labels.
- `-o <output_tree.dt>` → Path to save the trained model (JSON format).

**Example:**
```sh
./dt -c train -i datasets/train.csv -t class -o model.dt
```

### 2. Making Predictions

```sh
./dt -c predict -i <prediction_data_file.csv> -m <model_file.dt> -o <predictions.csv>
```

**Arguments:**
- `-c predict` → Specifies the prediction command.
- `-i <prediction_data_file.csv>` → Path to the dataset for predictions.
- `-m <model_file.dt>` → Path to the trained model file.
- `-o <predictions.csv>` → Path to save predictions.

**Example:**
```sh
./dt -c predict -i datasets/test.csv -m model.dt -o predictions.csv
```

## Input Requirements

- The dataset must be in **CSV format** with a header row.
- Feature columns may include **categorical, numeric, date, or timestamp** values.
- The **target column** must be specified during training.
- The trained model is saved in **JSON format**.
- The test dataset for predictions should have the **same feature columns** as the training dataset.

## Error Handling

| Error | Cause | Fix |
|-------|-------|-----|
| `Error: Missing input file` | Input CSV file is missing | Check the file path |
| `Error: Target column not found` | The specified column is not in the dataset | Verify the CSV column names |
| `Error: Model file not found` | The specified model file does not exist | Train a model first or check the file path |
| `Error: Output path not specified` | No output file provided | Specify an output file path |

## Testing

The project includes **unit tests** to validate correctness and performance.

Run the tests using:

```sh
go test ./...
```

## Authors

- Rodney Ochieng ([GitHub](https://github.com/rodneyo1))
- Valeria Muhembele ([GitHub](https://github.com/anamivale))
- Sheilla Juma ([GitHub](https://github.com/a-j-sheilla))
- Moses Onyango ([GitHub](https://github.com/moseeh))
- Thadeus Ogondola ([GitHub](https://github.com/TMassive42))

