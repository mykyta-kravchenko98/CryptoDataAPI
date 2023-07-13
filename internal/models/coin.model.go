package models

import (
	pb "github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/cryptodata_v1"
)

type GetAllCoinsResponse struct {
	Coins []CoinDto `json:"data"`
}

type CoinDto struct {
	Id     int64    `json:"id"`
	Name   string   `json:"name"`
	Symbol string   `json:"symbol"`
	Rank   uint32   `json:"cmc_rank"`
	Quotes QuoteDto `json:"quote"`
}

type QuoteDto struct {
	USD QuoteDataDto `json:"USD"`
}

type QuoteDataDto struct {
	Price            float64 `json:"price"`
	Volume24h        float64 `json:"volume_24h"`
	VolumeChange24h  float64 `json:"volume_change_24h"`
	PercentChange1h  float64 `json:"percent_change_1h"`
	PercentChange24h float64 `json:"percent_change_24h"`
	PercentChange7d  float64 `json:"percent_change_7d"`
	MarketCap        float64 `json:"market_cap"`
}

type Coin struct {
	Id               int64   `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	Rank             uint32  `json:"rank"`
	PriceUSD         float64 `json:"price"`
	Volume24h        float64 `json:"volume24h"`
	VolumeChange24h  float64 `json:"volumeChange24h"`
	PercentChange1h  float64 `json:"percentChange1h"`
	PercentChange24h float64 `json:"percentChange24h"`
	PercentChange7d  float64 `json:"percentChange7d"`
	MarketCap        float64 `json:"marketCap"`
}

func (c *Coin) GetDataFromDto(coinDto CoinDto) {
	c.Id = coinDto.Id
	c.Name = coinDto.Name
	c.Rank = coinDto.Rank
	c.Symbol = coinDto.Symbol
	c.PriceUSD = coinDto.Quotes.USD.Price
	c.Volume24h = coinDto.Quotes.USD.Volume24h
	c.VolumeChange24h = coinDto.Quotes.USD.VolumeChange24h
	c.PercentChange1h = coinDto.Quotes.USD.PercentChange1h
	c.PercentChange24h = coinDto.Quotes.USD.PercentChange24h
	c.PercentChange7d = coinDto.Quotes.USD.PercentChange7d
	c.MarketCap = coinDto.Quotes.USD.MarketCap
}

func (c *Coin) ProtoCoin() *pb.Coin {
	pc := pb.Coin{}

	pc.Id = c.Id
	pc.Name = c.Name
	pc.MarketCap = c.MarketCap
	pc.PercentChange1H = c.PercentChange1h
	pc.PercentChange24H = c.PercentChange24h
	pc.PercentChange7D = c.PercentChange7d
	pc.Price = c.PriceUSD
	pc.Rank = c.Rank
	pc.Symbol = c.Symbol
	pc.Volume24H = c.Volume24h
	pc.VolumeChange24H = c.VolumeChange24h

	return &pc
}
