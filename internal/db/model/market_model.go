package model

type Market struct {
	ID            int64  `db:"id"`
	Question      string `db:"question"`
	Outcomes      string `db:"outcomes"`       // e.g. "[\"Yes\",\"No\"]"
	OutcomePrices string `db:"outcome_prices"` // e.g. "[\"0.115\",\"0.885\"]"
	Volume        string `db:"volume"`         // e.g. "19.8129"
	ImageURL      string `db:"image_url"`
}
