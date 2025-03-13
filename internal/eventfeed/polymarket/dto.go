package polymarket

type Event struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	StartDate   string   `json:"startDate"`
	Image       string   `json:"image"`
	Volume      float64  `json:"volume"`
	Tags        []Tag    `json:"tags"`
	Markets     []Market `json:"markets"`
}

type Tag struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Market struct {
	ID            string  `json:"id"`
	Question      string  `json:"question"`
	Description   string  `json:"description"`
	Outcomes      string  `json:"outcomes"` // This might be a JSON string like "[\"Yes\",\"No\"]"
	OutcomePrices string  `json:"outcomePrices"`
	Volume        string  `json:"volume"`
	VolumeNum     float64 `json:"volumeNum"`
	Featured      bool    `json:"featured"`
	Active        bool    `json:"active"`
	Closed        bool    `json:"closed"`
}
