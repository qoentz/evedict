package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/qoentz/evedict/internal/db/model"
)

type ForecastRepository struct {
	DB *sqlx.DB
}

func NewForecastRepository(db *sqlx.DB) *ForecastRepository {
	return &ForecastRepository{
		DB: db,
	}
}

func (r *ForecastRepository) GetForecasts(limit int, offset int) ([]model.Forecast, error) {
	var forecasts []model.Forecast
	query := `
		SELECT id, headline, summary, image_url, timestamp 
		FROM forecast
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`

	err := r.DB.Select(&forecasts, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecasts: %v", err)
	}

	return forecasts, nil
}

func (r *ForecastRepository) GetForecast(forecastId uuid.UUID) (*model.Forecast, error) {
	var forecast model.Forecast
	query := `
		SELECT id, headline, summary, image_url, timestamp 
		FROM forecast
		WHERE id = $1
	`
	err := r.DB.Get(&forecast, query, forecastId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %v", err)
	}

	outcomes, err := r.getOutcomesByForecastID(forecastId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch outcomes for forecast %s: %v", forecastId, err)
	}
	forecast.Outcomes = outcomes

	sources, err := r.getSourcesByForecastID(forecastId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources for forecast %s: %v", forecastId, err)
	}
	forecast.Sources = sources

	return &forecast, nil
}

func (r *ForecastRepository) SaveForecast(forecast *model.Forecast) error {
	// Start a transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Insert the main Forecast record with a specified UUID
	query := `INSERT INTO forecast (id, headline, summary, image_url, timestamp) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, forecast.ID, forecast.Headline, forecast.Summary, forecast.ImageURL, forecast.Timestamp)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert forecast: %v", err)
	}

	// Insert associated Outcomes with specified UUIDs
	outcomeQuery := `INSERT INTO outcome (id, forecast_id, content, confidence_level) VALUES ($1, $2, $3, $4)`
	for _, outcome := range forecast.Outcomes {
		_, err = tx.Exec(outcomeQuery, outcome.ID, forecast.ID, outcome.Content, outcome.ConfidenceLevel)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert outcome: %v", err)
		}
	}

	// Insert associated Sources with specified UUIDs
	sourceQuery := `INSERT INTO source (id, forecast_id, name, title, url) VALUES ($1, $2, $3, $4, $5)`
	for _, source := range forecast.Sources {
		_, err = tx.Exec(sourceQuery, source.ID, forecast.ID, source.Name, source.Title, source.URL)
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

func (r *ForecastRepository) SaveForecasts(forecasts []model.Forecast) error {
	// Start a transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Prepare queries, including the `id` field
	forecastQuery := `INSERT INTO forecast (id, headline, summary, image_url, timestamp) VALUES ($1, $2, $3, $4, $5)`
	outcomeQuery := `INSERT INTO outcome (id, forecast_id, content, confidence_level) VALUES ($1, $2, $3, $4)`
	sourceQuery := `INSERT INTO source (id, forecast_id, name, title, url, image_url) VALUES ($1, $2, $3, $4, $5, $6)`

	// Insert each forecast and associated records within the transaction
	for i := range forecasts {
		forecast := &forecasts[i]

		// Insert the main Forecast record with a specified UUID
		_, err = tx.Exec(forecastQuery, forecast.ID, forecast.Headline, forecast.Summary, forecast.ImageURL, forecast.Timestamp)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert forecast: %v", err)
		}

		// Insert associated Outcomes with specified UUIDs
		for _, outcome := range forecast.Outcomes {
			_, err = tx.Exec(outcomeQuery, outcome.ID, forecast.ID, outcome.Content, outcome.ConfidenceLevel)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert outcome: %v", err)
			}
		}

		// Insert associated Sources with specified UUIDs
		for _, source := range forecast.Sources {
			_, err = tx.Exec(sourceQuery, source.ID, forecast.ID, source.Name, source.Title, source.URL, source.ImageURL)
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

func (r *ForecastRepository) getOutcomesByForecastID(forecastID uuid.UUID) ([]model.Outcome, error) {
	var outcomes []model.Outcome
	err := r.DB.Select(&outcomes, `SELECT id, forecast_id, content, confidence_level FROM outcome WHERE forecast_id = $1`, forecastID)
	return outcomes, err
}

func (r *ForecastRepository) getSourcesByForecastID(forecastID uuid.UUID) ([]model.Source, error) {
	var sources []model.Source
	err := r.DB.Select(&sources, `SELECT id, forecast_id, name, title, url FROM source WHERE forecast_id = $1`, forecastID)
	return sources, err
}
