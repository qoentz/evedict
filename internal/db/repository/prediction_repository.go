package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/qoentz/evedict/internal/db/model"
)

type PredictionRepository struct {
	DB *sqlx.DB
}

func NewPredictionRepository(db *sqlx.DB) *PredictionRepository {
	return &PredictionRepository{
		DB: db,
	}
}

func (r *PredictionRepository) GetPredictions() ([]model.Prediction, error) {
	var predictions []model.Prediction
	err := r.DB.Select(&predictions, `SELECT id, headline, summary, image_url, timestamp FROM prediction`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch predictions: %v", err)
	}

	for i := range predictions {
		outcomes, err := r.getOutcomesByPredictionID(predictions[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch outcomes for prediction %d: %v", predictions[i].ID, err)
		}
		predictions[i].Outcomes = outcomes

		sources, err := r.getSourcesByPredictionID(predictions[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch sources for prediction %d: %v", predictions[i].ID, err)
		}
		predictions[i].Sources = sources
	}

	return predictions, nil
}

func (r *PredictionRepository) SavePrediction(prediction *model.Prediction) error {
	// Start a transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Insert the main Prediction record with a specified UUID
	query := `INSERT INTO prediction (id, headline, summary, image_url, timestamp) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, prediction.ID, prediction.Headline, prediction.Summary, prediction.ImageURL, prediction.Timestamp)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert prediction: %v", err)
	}

	// Insert associated Outcomes with specified UUIDs
	outcomeQuery := `INSERT INTO outcome (id, prediction_id, content, confidence_level) VALUES ($1, $2, $3, $4)`
	for _, outcome := range prediction.Outcomes {
		_, err = tx.Exec(outcomeQuery, outcome.ID, prediction.ID, outcome.Content, outcome.ConfidenceLevel)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert outcome: %v", err)
		}
	}

	// Insert associated Sources with specified UUIDs
	sourceQuery := `INSERT INTO source (id, prediction_id, name, title, url) VALUES ($1, $2, $3, $4, $5)`
	for _, source := range prediction.Sources {
		_, err = tx.Exec(sourceQuery, source.ID, prediction.ID, source.Name, source.Title, source.URL)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert source: %v", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *PredictionRepository) SavePredictions(predictions []model.Prediction) error {
	// Start a transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Prepare queries, including the `id` field
	predictionQuery := `INSERT INTO prediction (id, headline, summary, image_url, timestamp) VALUES ($1, $2, $3, $4, $5)`
	outcomeQuery := `INSERT INTO outcome (id, prediction_id, content, confidence_level) VALUES ($1, $2, $3, $4)`
	sourceQuery := `INSERT INTO source (id, prediction_id, name, title, url) VALUES ($1, $2, $3, $4, $5)`

	// Insert each prediction and associated records within the transaction
	for i := range predictions {
		prediction := &predictions[i]

		// Insert the main Prediction record with a specified UUID
		_, err = tx.Exec(predictionQuery, prediction.ID, prediction.Headline, prediction.Summary, prediction.ImageURL, prediction.Timestamp)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert prediction: %v", err)
		}

		// Insert associated Outcomes with specified UUIDs
		for _, outcome := range prediction.Outcomes {
			_, err = tx.Exec(outcomeQuery, outcome.ID, prediction.ID, outcome.Content, outcome.ConfidenceLevel)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert outcome: %v", err)
			}
		}

		// Insert associated Sources with specified UUIDs
		for _, source := range prediction.Sources {
			_, err = tx.Exec(sourceQuery, source.ID, prediction.ID, source.Name, source.Title, source.URL)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert source: %v", err)
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *PredictionRepository) getOutcomesByPredictionID(predictionID uuid.UUID) ([]model.Outcome, error) {
	var outcomes []model.Outcome
	err := r.DB.Select(&outcomes, `SELECT id, prediction_id, content, confidence_level FROM outcome WHERE prediction_id = $1`, predictionID)
	return outcomes, err
}

func (r *PredictionRepository) getSourcesByPredictionID(predictionID uuid.UUID) ([]model.Source, error) {
	var sources []model.Source
	err := r.DB.Select(&sources, `SELECT id, prediction_id, name, title, url FROM source WHERE prediction_id = $1`, predictionID)
	return sources, err
}
