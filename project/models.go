package project

type Currency struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Rate     float64 `json:"rate"`
	InsertDt string  `json:"insert_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
