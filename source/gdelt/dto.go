package gdelt

type Article struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	SeenDate string `json:"seendate"`
	Domain   string `json:"domain"`
	Language string `json:"language"`
}

type Response struct {
	Articles []Article `json:"articles"`
}
