package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/configs"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/models"
	lrucache "github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/cache"
	"github.com/mykyta-kravchenko98/ValueShift/pkg/clients/rest"
)

type syncService struct {
	cache  *lrucache.LRUCache
	config configs.CoinMarketCapConfig
}

type SyncService interface {
	SyncTop50CoinMarketCurrency() error
}

func NewSyncService(c *lrucache.LRUCache, conf configs.CoinMarketCapConfig) SyncService {
	instance := &syncService{
		cache:  c,
		config: conf,
	}

	return instance
}

func (s *syncService) SyncTop50CoinMarketCurrency() error {
	headers := http.Header{}
	headers.Set("X-CMC_PRO_API_KEY", s.config.APIKey)
	headers.Set("Accept", "application/json")

	url := fmt.Sprintf("%s/cryptocurrency/listings/latest?limit=50", s.config.URL)

	resp, err := rest.Get(url, headers)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response models.GetAllCoinsResponse
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	for _, coinDto := range response.Coins {
		newCoin := models.Coin{}
		newCoin.GetDataFromDto(coinDto)
		s.cache.Put(int(newCoin.Rank), newCoin)
	}

	return nil
}
