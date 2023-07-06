package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/configs"
	grpcserver "github.com/mykyta-kravchenko98/CryptoDataAPI/internal/grpc_server"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/services"
	lrucache "github.com/mykyta-kravchenko98/CryptoDataAPI/pkg/cache"
	"github.com/rs/zerolog/log"
)

func main() {
	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	//Load configuration
	conf, confErr := configs.LoadConfigs(env)
	if confErr != nil {
		log.Fatal().Err(confErr).Msg("Config load failed")
	}

	cache := lrucache.InitLRUCache(50)
	syncService := services.NewSyncService(&cache, conf.CoinMarketCap)
	dataService := services.NewDataService(&cache)

	//Init Cache
	if len(cache.Cache) <= 0 {
		syncService.SyncTop50CoinMarketCurrency()
	}

	//Sync Job
	job := gocron.NewScheduler(time.UTC)
	job.Every(5).Minute().Do(func() {
		syncService.SyncTop50CoinMarketCurrency()
	})

	job.StartAsync()

	//Init gRPC server
	if err := grpcserver.Init(dataService, conf.Server); err != nil {
		log.Fatal().Err(err)
	}
}
