package services

import (
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/models"
	lrucache "github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/cache"
	pb "github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/cryptodata_v1"
)

type dataService struct {
	cache *lrucache.LRUCache
}

type DataService interface {
	GetTop50CoinMarketCurrency() ([]*pb.Coin, error)
}

func NewDataService(c *lrucache.LRUCache) DataService {
	instance := &dataService{
		cache: c,
	}

	return instance
}

func (ds *dataService) GetTop50CoinMarketCurrency() ([]*pb.Coin, error) {
	coins := make([]*pb.Coin, 0, 50)
	for i := 1; i <= 50; i++ {
		if coin, ok := ds.cache.Get(i).(models.Coin); ok {
			pc := coin.ProtoCoin()
			coins = append(coins, pc)
		}
	}

	return coins, nil
}
