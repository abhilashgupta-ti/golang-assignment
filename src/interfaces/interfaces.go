package interfaces

// Structs to construct the JSON API response
type BitcoinRatesStruct struct {
	Eur string `json:"EUR"`
	Usd string `json:"USD"`
}

type DataStruct struct {
	Bitcoin BitcoinRatesStruct `json:"bitcoin"`
}

type ResponseStruct struct {
	Data DataStruct `json:"data"`
}

type PriceTracker interface {
	GetPriceFromTracker() (ResponseStruct, error)
	GetName() string
	GetURL() string
}
