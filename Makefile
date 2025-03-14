.PHONY: build train predict clean all

# Variables
BINARY=dt
TRAIN_DATA=datasets/pixar_films new.csv
MODEL_FILE=model.dt
PREDICT_OUTPUT=predictions.csv
TARGET_COLUMN=box_office_worldwide

build:
	@echo "Building decision tree with go build."
	@go build -o $(BINARY)

train: build
	@echo "Training model with $(TRAIN_DATA)..."
	./$(BINARY) -c train -i "$(TRAIN_DATA)" -t $(TARGET_COLUMN) -o $(MODEL_FILE)

predict: build
	@echo "Making predictions..."
	./$(BINARY) -c predict -i "$(TRAIN_DATA)" -m $(MODEL_FILE) -o $(PREDICT_OUTPUT)

clean:
	@echo "Cleaning up files..."
	@rm -f $(BINARY) $(MODEL_FILE) $(PREDICT_OUTPUT)

all: clean build train predict

# Optional: Add a 'run' target for quick testing
run: predict
	@echo "Predictions written to $(PREDICT_OUTPUT)"