package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/qoentz/evedict/internal/db/model"
	"github.com/qoentz/evedict/internal/util"
)

type ForecastRepository struct {
	DB *sqlx.DB
}

func NewForecastRepository(db *sqlx.DB) *ForecastRepository {
	return &ForecastRepository{
		DB: db,
	}
}

func (r *ForecastRepository) GetForecasts(limit int, offset int, category *util.Category, isApproved bool) ([]model.Forecast, error) {
	var forecasts []model.Forecast
	var err error

	if category != nil {
		query := `
			SELECT id, headline, summary, image_url, timestamp 
			FROM forecast
			WHERE category = $1
			AND is_approved = $4
			ORDER BY timestamp DESC
			LIMIT $2 OFFSET $3
		`
		err = r.DB.Select(&forecasts, query, *category, limit, offset, isApproved)
	} else {
		query := `
			SELECT id, headline, summary, image_url, timestamp 
			FROM forecast
			WHERE is_approved = $3
			ORDER BY timestamp DESC
			LIMIT $1 OFFSET $2
		`
		err = r.DB.Select(&forecasts, query, limit, offset, isApproved)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecasts: %v", err)
	}

	return forecasts, nil
}

func (r *ForecastRepository) GetForecast(forecastID uuid.UUID) (*model.Forecast, error) {
	var f model.Forecast
	forecastQuery := `
        SELECT id, headline, summary, image_url, category, timestamp
        FROM forecast
        WHERE id = $1
    `
	err := r.DB.Get(&f, forecastQuery, forecastID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast: %v", err)
	}

	outcomes, err := r.getOutcomesByForecastID(forecastID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch outcomes: %v", err)
	}
	f.Outcomes = outcomes

	sources, err := r.getSourcesByForecastID(forecastID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sources: %v", err)
	}
	f.Sources = sources

	tags, err := r.getTagsByForecastID(forecastID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags for forecast %s: %v", forecastID, err)
	}
	f.Tags = tags

	m, err := r.getMarketByForecastID(forecastID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch market: %v", err)
	}
	f.Market = m

	return &f, nil
}

func (r *ForecastRepository) GetRelatedForecastsByTagAndCategory(
	mainID uuid.UUID,
	tagNames []string,
	category util.Category,
	limit int,
) ([]model.RelatedForecast, error) {

	// If you’re using lib/pq, you can pass tagNames as pq.StringArray to "t2.name = ANY($2)"
	// We'll define matched_by_tag with a boolean literal in each SELECT.
	// The final ORDER BY puts matched_by_tag = TRUE first, then newest to oldest.

	unionQuery := `SELECT *
FROM (
    SELECT DISTINCT ON (id) id, headline, summary, image_url, timestamp, matched_by_tag
    FROM (
        -- Tag-matched forecasts
        SELECT f.id, f.headline, f.summary, f.image_url, f.timestamp, TRUE AS matched_by_tag
        FROM forecast f
        JOIN forecast_tag ft ON ft.forecast_id = f.id
        JOIN tag t2 ON t2.id = ft.tag_id
        WHERE f.id <> $1
          AND t2.name = ANY($2)
    		AND is_approved = true

        UNION ALL

        -- Category-matched forecasts
        SELECT f.id, f.headline, f.summary, f.image_url, f.timestamp, FALSE AS matched_by_tag
        FROM forecast f
        WHERE f.id <> $1
          AND f.category = $3
          AND is_approved = true
    ) AS unioned
    ORDER BY id, matched_by_tag DESC, timestamp DESC
) AS deduped
ORDER BY matched_by_tag DESC, timestamp DESC
LIMIT $4;
`

	var related []model.RelatedForecast
	err := r.DB.Select(&related, unionQuery,
		mainID,                   // $1
		pq.StringArray(tagNames), // $2 (the list of tags)
		category,                 // $3
		limit,                    // $4
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch related forecasts: %v", err)
	}

	return related, nil
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

func (r *ForecastRepository) SavePolyForecasts(forecasts []model.Forecast) error {
	// Start transaction
	tx, err := r.DB.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Forecast INSERT query
	forecastQuery := `
        INSERT INTO forecast (id, headline, summary, image_url, category, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	outcomeQuery := `
        INSERT INTO outcome (id, forecast_id, content, confidence_level)
        VALUES ($1, $2, $3, $4)
    `

	sourceQuery := `
        INSERT INTO source (id, forecast_id, name, title, url, image_url)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	tagUpsertQuery := `
        INSERT INTO tag (name)
        VALUES ($1)
        ON CONFLICT (name)
        DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `
	forecastTagQuery := `
        INSERT INTO forecast_tag (forecast_id, tag_id)
        VALUES ($1, $2)
    `

	// Market INSERT query
	marketQuery := `
    INSERT INTO market (id, question, outcomes, outcome_prices, volume, image_url)
    VALUES ($1, $2, $3, $4, $5, $6)
    ON CONFLICT (id) 
    DO UPDATE SET 
        question = EXCLUDED.question,
        outcomes = EXCLUDED.outcomes,
        outcome_prices = EXCLUDED.outcome_prices,
        volume = EXCLUDED.volume,
        image_url = EXCLUDED.image_url;
`

	// Forecast-Market relation insert query
	forecastMarketQuery := `
        INSERT INTO forecast_market (forecast_id, market_id)
        VALUES ($1, $2)
    `

	// Insert each forecast + associated records
	for i := range forecasts {
		forecast := &forecasts[i]

		// === INSERT MAIN FORECAST ===
		_, err = tx.Exec(forecastQuery,
			forecast.ID,
			forecast.Headline,
			forecast.Summary,
			forecast.ImageURL,
			forecast.Category,
			forecast.Timestamp,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert forecast: %v", err)
		}

		// === OUTCOMES ===
		for _, outcome := range forecast.Outcomes {
			_, err = tx.Exec(outcomeQuery,
				outcome.ID,
				forecast.ID,
				outcome.Content,
				outcome.ConfidenceLevel,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert outcome: %v", err)
			}
		}

		// === TAGS ===
		for _, tag := range forecast.Tags {
			// 1. Upsert the tag by name
			var tagID uuid.UUID
			err = tx.QueryRow(tagUpsertQuery, tag.Name).Scan(&tagID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to upsert tag (name=%q): %v", tag.Name, err)
			}

			// 2. Insert into forecast_tag
			_, err = tx.Exec(forecastTagQuery, forecast.ID, tagID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert forecast_tag link: %v", err)
			}
		}

		// === SOURCES ===
		for _, source := range forecast.Sources {
			_, err = tx.Exec(sourceQuery,
				source.ID,
				forecast.ID,
				source.Name,
				source.Title,
				source.URL,
				source.ImageURL,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert source: %v", err)
			}
		}

		// === MARKET (1:1) ===
		if forecast.Market != nil {
			// Insert into market table
			_, err = tx.Exec(marketQuery,
				forecast.Market.ID,
				forecast.Market.Question,
				forecast.Market.Outcomes,
				forecast.Market.OutcomePrices,
				forecast.Market.Volume,
				forecast.Market.ImageURL,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert market: %v", err)
			}

			// Insert into forecast_market join table
			_, err = tx.Exec(forecastMarketQuery, forecast.ID, forecast.Market.ID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert into forecast_market relation: %v", err)
			}
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

	// 1) Prepare the forecast INSERT query (note the "category" field is included now)
	forecastQuery := `
        INSERT INTO forecast (id, headline, summary, image_url, category, timestamp)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	// 2) Prepare the others (same as before)
	outcomeQuery := `
        INSERT INTO outcome (id, forecast_id, content, confidence_level)
        VALUES ($1, $2, $3, $4)
    `
	sourceQuery := `
        INSERT INTO source (id, forecast_id, name, title, url, image_url)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	// 3) We'll use this to "upsert" tags by name
	tagUpsertQuery := `
        INSERT INTO tag (name)
        VALUES ($1)
        ON CONFLICT (name)
        DO UPDATE SET name = EXCLUDED.name
        RETURNING id
    `
	forecastTagQuery := `
        INSERT INTO forecast_tag (forecast_id, tag_id)
        VALUES ($1, $2)
    `

	// Insert each forecast + associated records
	for i := range forecasts {
		forecast := &forecasts[i]

		// === INSERT MAIN FORECAST (with category) ===
		_, err = tx.Exec(forecastQuery,
			forecast.ID,
			forecast.Headline,
			forecast.Summary,
			forecast.ImageURL,
			forecast.Category, // <--- category now included here
			forecast.Timestamp,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert forecast: %v", err)
		}

		// === OUTCOMES ===
		for _, outcome := range forecast.Outcomes {
			_, err = tx.Exec(outcomeQuery,
				outcome.ID,
				forecast.ID,
				outcome.Content,
				outcome.ConfidenceLevel,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert outcome: %v", err)
			}
		}

		// === TAGS ===
		for _, tag := range forecast.Tags {
			// 1. Upsert the tag by name, let Postgres assign/keep its UUID
			var tagID uuid.UUID
			err = tx.QueryRow(tagUpsertQuery, tag.Name).Scan(&tagID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to upsert tag (name=%q): %v", tag.Name, err)
			}

			// 2. Insert into forecast_tag
			_, err = tx.Exec(forecastTagQuery, forecast.ID, tagID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to insert forecast_tag link: %v", err)
			}
		}

		// === SOURCES ===
		for _, source := range forecast.Sources {
			_, err = tx.Exec(sourceQuery,
				source.ID,
				forecast.ID,
				source.Name,
				source.Title,
				source.URL,
				source.ImageURL,
			)
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

func (r *ForecastRepository) MarkForecastApproved(forecastID uuid.UUID) error {
	_, err := r.DB.Exec(`
		UPDATE forecast
		SET is_approved = true
		WHERE id = $1
	`, forecastID)
	return err
}

func (r *ForecastRepository) getOutcomesByForecastID(forecastID uuid.UUID) ([]model.Outcome, error) {
	var outcomes []model.Outcome
	err := r.DB.Select(&outcomes, `SELECT id, forecast_id, content, confidence_level FROM outcome WHERE forecast_id = $1`, forecastID)
	return outcomes, err
}

func (r *ForecastRepository) getTagsByForecastID(forecastID uuid.UUID) ([]model.Tag, error) {
	var tags []model.Tag
	query := `
        SELECT t.id, t.name
        FROM tag t
        JOIN forecast_tag ft ON ft.tag_id = t.id
        WHERE ft.forecast_id = $1
    `
	err := r.DB.Select(&tags, query, forecastID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %v", err)
	}
	return tags, nil
}

func (r *ForecastRepository) getSourcesByForecastID(forecastID uuid.UUID) ([]model.Source, error) {
	var sources []model.Source
	err := r.DB.Select(&sources, `SELECT id, forecast_id, name, title, url, image_url FROM source WHERE forecast_id = $1`, forecastID)
	return sources, err
}

func (r *ForecastRepository) getMarketByForecastID(forecastID uuid.UUID) (*model.Market, error) {
	var marketID int64

	// Step 1: Get the market_id from forecast_market
	getMarketIDQuery := `
        SELECT market_id FROM forecast_market WHERE forecast_id = $1
    `
	err := r.DB.Get(&marketID, getMarketIDQuery, forecastID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No market linked to this forecast
		}
		return nil, fmt.Errorf("failed to fetch market_id: %v", err)
	}

	// Step 2: Fetch the market details using the retrieved market_id
	var market model.Market
	getMarketQuery := `
        SELECT id, question, outcomes, outcome_prices, volume, image_url
        FROM market
        WHERE id = $1
    `
	err = r.DB.Get(&market, getMarketQuery, marketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Market does not exist, return nil
		}
		return nil, fmt.Errorf("failed to fetch market details: %v", err)
	}

	return &market, nil
}

func (r *ForecastRepository) CheckImageURL(imageURL string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM forecast WHERE image_url = $1 AND is_approved = TRUE)`
	err := r.DB.QueryRow(query, imageURL).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}
