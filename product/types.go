package product

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type HitsTotal struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type SearchHits struct {
	Total HitsTotal `json:"total"`
}
