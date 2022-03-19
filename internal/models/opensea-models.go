package models

type OpenSeaStats struct {
	Stats `json:"stats"`
}

type Stats struct {
	OneDayVolume          float64 `json:"one_day_volume"`
	OneDayChange          float64 `json:"one_day_change"`
	OneDaySales           float64 `json:"one_day_sales"`
	OneDayAveragePrice    float64 `json:"one_day_average_price"`
	SevenDayVolume        float64 `json:"seven_day_volume"`
	SevenDayChange        float64 `json:"seven_day_change"`
	SevenDaySales         float64 `json:"seven_day_sales"`
	SevenDayAveragePrice  float64 `json:"seven_day_average_price"`
	ThirtyDayVolume       float64 `json:"thirty_day_volume"`
	ThirtyDayChange       float64 `json:"thirty_day_change"`
	ThirtyDaySales        float64 `json:"thirty_day_sales"`
	ThirtyDayAveragePrice float64 `json:"thirty_day_average_price"`
	TotalVolume           float64 `json:"total_volume"`
	TotalSales            float64 `json:"total_sales"`
	TotalSupply           float64 `json:"total_supply"`
	Count                 float64 `json:"count"`
	NumOwners             int     `json:"num_owners"`
	AveragePrice          float64 `json:"average_price"`
	NumReports            int     `json:"num_reports"`
	MarketCap             float64 `json:"market_cap"`
	FloorPrice            float64 `json:"floor_price"`
}
