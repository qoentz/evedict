package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/qoentz/evedict/internal/db/model"
)

type DivinationRepository struct {
	DB *sqlx.DB
}

func NewDivinationRepository(db *sqlx.DB) *DivinationRepository {
	return &DivinationRepository{
		DB: db,
	}
}

func (r *DivinationRepository) GetDivinations() ([]model.Divination, error) {
	var divinations []model.Divination
	err := r.DB.Select(&divinations, `
		SELECT id, headline, summary, image_url, timestamp 
		FROM divination
		ORDER BY timestamp DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch divinations: %v", err)
	}

	for i := range divinations {
		outcomes, err := r.getOutcomesByDivinationID(divinations[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch outcomes for divination %d: %v", divinations[i].ID, err)
		}
		divinations[i].Outcomes = outcomes

		sources, err := r.getSourcesByDivinationID(divinations[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch sources for divination %d: %v", divinations[i].ID, err)
		}
		divinations[i].Sources = sources
	}

	return divinations, nil
}

func (r *DivinationRepository) SaveDivination(divination *model.Divination) error {
	// Start a transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Insert the main Divination record with a specified UUID
	query := `INSERT INTO divination (id, headline, summary, image_url, timestamp) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(query, divination.ID, divination.Headline, divination.Summary, divination.ImageURL, divination.Timestamp)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert divination: %v", err)
	}

	// Insert associated Outcomes with specified UUIDs
	outcomeQuery := `INSERT INTO outcome (id, divination_id, content, confidence_level) VALUES ($1, $2, $3, $4)`
	for _, outcome := range divination.Outcomes {
		_, err = tx.Exec(outcomeQuery, outcome.ID, divination.ID, outcome.Content, outcome.ConfidenceLevel)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert outcome: %v", err)
		}
	}

	// Insert associated Sources with specified UUIDs
	sourceQuery := `INSERT INTO source (id, divination_id, name, title, url) VALUES ($1, $2, $3, $4, $5)`
	for _, source := range divination.Sources {
		_, err = tx.Exec(sourceQuery, source.ID, divination.ID, source.Name, source.Title, source.URL)
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

func (r *DivinationRepository) SaveDivinations(divinations []model.Divination) error {
	// Start a transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Prepare queries, including the `id` field
	divinationQuery := `INSERT INTO divination (id, headline, summary, image_url, timestamp) VALUES ($1, $2, $3, $4, $5)`
	outcomeQuery := `INSERT INTO outcome (id, divination_id, content, confidence_level) VALUES ($1, $2, $3, $4)`
	sourceQuery := `INSERT INTO source (id, divination_id, name, title, url) VALUES ($1, $2, $3, $4, $5)`

	// Insert each divination and associated records within the transaction
	for i := range divinations {
		divination := &divinations[i]

		// Insert the main Divination record with a specified UUID
		_, err = tx.Exec(divinationQuery, divination.ID, divination.Headline, divination.Summary, divination.ImageURL, divination.Timestamp)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert divination: %v", err)
		}

		// Insert associated Outcomes with specified UUIDs
		for _, outcome := range divination.Outcomes {
			_, err = tx.Exec(outcomeQuery, outcome.ID, divination.ID, outcome.Content, outcome.ConfidenceLevel)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert outcome: %v", err)
			}
		}

		// Insert associated Sources with specified UUIDs
		for _, source := range divination.Sources {
			_, err = tx.Exec(sourceQuery, source.ID, divination.ID, source.Name, source.Title, source.URL)
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

func (r *DivinationRepository) getOutcomesByDivinationID(divinationID uuid.UUID) ([]model.Outcome, error) {
	var outcomes []model.Outcome
	err := r.DB.Select(&outcomes, `SELECT id, divination_id, content, confidence_level FROM outcome WHERE divination_id = $1`, divinationID)
	return outcomes, err
}

func (r *DivinationRepository) getSourcesByDivinationID(divinationID uuid.UUID) ([]model.Source, error) {
	var sources []model.Source
	err := r.DB.Select(&sources, `SELECT id, divination_id, name, title, url FROM source WHERE divination_id = $1`, divinationID)
	return sources, err
}
