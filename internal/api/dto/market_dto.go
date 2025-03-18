package dto

type Market struct {
	ID                int64    `json:"id"`
	Question          string   `json:"question"`
	Outcomes          string   `json:"outcomes"`
	OutcomePrices     string   `json:"outcomePrices"`
	Volume            string   `json:"volume"`
	ImageURL          string   `json:"imageUrl"`
	OutcomeList       []string `json:"-"`
	OutcomePricesList []string `json:"-"`
}
