package entity

type FetchOneResponse struct {
	Base   string             `json:"base"`
	Result map[string]float64 `json:"result"` // example: results: {"CNY": 7.34}
}

type FetchMultiResponse struct {
	Base    string             `json:"base"`
	Results map[string]float64 `json:"results"` // example: results: {"CNY": 7.34}
}

type CryptoFetchPrices struct {
	Prices map[string]float64 `json:"prices"`
}
