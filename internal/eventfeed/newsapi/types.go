package newsapi

type Endpoint string

const (
	TopHeadlines Endpoint = "top-headlines"
	Everything   Endpoint = "everything"
)

type Category string

const (
	Business      Category = "business"
	Entertainment Category = "entertainment"
	General       Category = "general"
	Health        Category = "health"
	Science       Category = "science"
	Sports        Category = "sports"
	Technology    Category = "technology"
)
