package models

type ReservoirCollectionResponse struct {
	Collection struct {
		ID       string `json:"id"`
		Slug     string `json:"slug"`
		Name     string `json:"name"`
		Metadata struct {
			ImageURL        string `json:"imageUrl"`
			DiscordURL      string `json:"discordUrl"`
			Description     string `json:"description"`
			ExternalURL     string `json:"externalUrl"`
			BannerImageURL  string `json:"bannerImageUrl"`
			TwitterUsername string `json:"twitterUsername"`
		} `json:"metadata"`
		SampleImages []interface{} `json:"sampleImages"`
		TokenCount   string        `json:"tokenCount"`
		OnSaleCount  string        `json:"onSaleCount"`
		TokenSetID   string        `json:"tokenSetId"`
		Royalties    struct {
			Bps       int    `json:"bps"`
			Recipient string `json:"recipient"`
		} `json:"royalties"`
		LastBuy struct {
			Value     interface{} `json:"value"`
			Timestamp interface{} `json:"timestamp"`
		} `json:"lastBuy"`
		LastSell struct {
			Value     float64 `json:"value"`
			Timestamp int     `json:"timestamp"`
		} `json:"lastSell"`
		FloorAsk struct {
			ID         string  `json:"id"`
			Price      float64 `json:"price"`
			Maker      string  `json:"maker"`
			ValidFrom  int     `json:"validFrom"`
			ValidUntil int     `json:"validUntil"`
			Token      struct {
				Contract string      `json:"contract"`
				TokenID  string      `json:"tokenId"`
				Name     interface{} `json:"name"`
				Image    interface{} `json:"image"`
			} `json:"token"`
		} `json:"floorAsk"`
		TopBid struct {
			ID         interface{} `json:"id"`
			Value      interface{} `json:"value"`
			Maker      interface{} `json:"maker"`
			ValidFrom  interface{} `json:"validFrom"`
			ValidUntil interface{} `json:"validUntil"`
		} `json:"topBid"`
		Rank struct {
			OneDay    int `json:"1day"`
			SevenDay  int `json:"7day"`
			Three0Day int `json:"30day"`
			AllTime   int `json:"allTime"`
		} `json:"rank"`
		Volume struct {
			OneDay    float64 `json:"1day"`
			SevenDay  float64 `json:"7day"`
			Three0Day float64 `json:"30day"`
			AllTime   float64 `json:"allTime"`
		} `json:"volume"`
	} `json:"collection"`
}

type ReservoirPriceMap struct {
	Tokens map[string]float64
}
