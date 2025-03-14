package algorithm

import (
	"sync"

	"decision-tree/models"
)

// Predict makes predictions for all records in the dataset
func Predict(tree *models.TreeNode) []interface{} {
	predictions := make([]interface{}, len(models.Records))

	// Use goroutines for parallel prediction
	var wg sync.WaitGroup
	workers := 4 // Number of worker goroutines
	batchSize := (len(models.Records) + workers - 1) / workers

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			start := workerID * batchSize
			end := (workerID + 1) * batchSize
			if end > len(models.Records) {
				end = len(models.Records)
			}

			for i := start; i < end; i++ {
				predictions[i] = predictRecord(models.Records[i], tree)
			}
		}(w)
	}

	wg.Wait()
	return predictions
}
