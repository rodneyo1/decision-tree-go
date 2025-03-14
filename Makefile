.PHONY: build train predict clean all

# Variables
BINARY=dt
TRAIN_DATA=datasets/loan_approval.csv
MODEL_FILE=model.dt
PREDICT_OUTPUT=predictions.csv
TARGET_COLUMN="Loan Approved"
PREDICT_DATA = datasets/sample.csv
build:
	@echo "Building decision tree with go build."
	@go build -o $(BINARY)

train: build
	@echo "Training model with $(TRAIN_DATA)..."
	./$(BINARY) -c train -i "$(TRAIN_DATA)" -t $(TARGET_COLUMN) -o $(MODEL_FILE)

predict: build
	@echo "Making predictions..."
	./$(BINARY) -c predict -i "$(PREDICT_DATA)" -m "$(MODEL_FILE)" -o $(PREDICT_OUTPUT)

clean:
	@echo "Cleaning up files..."
	@rm -f $(BINARY) $(MODEL_FILE) $(PREDICT_OUTPUT)

all: clean build train predict

# Optional: Add a 'run' target for quick testing
run: predict
	@echo "Predictions written to $(PREDICT_OUTPUT)"