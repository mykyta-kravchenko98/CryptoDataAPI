package main

import (
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/configs"
	grpcserver "github.com/mykyta-kravchenko98/CryptoDataAPI/internal/grpc_server"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/services"
	"github.com/mykyta-kravchenko98/CryptoDataAPI/internal/websocket"
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

	//create websocket and return websocket interface
	ws := websocket.StartWebSocket(dataService)

	//Sync Job
	job := gocron.NewScheduler(time.UTC)
	job.Every(5).Minute().Do(func() {
		coins, err := syncService.SyncTop50CoinMarketCurrency()
		if err == nil && ws.HasConnectedClients() {
			ws.SendMessage(websocket.CryptoCoinResponse{Coins: coins})
		}
	})

	job.StartAsync()

	//Init gRPC server
	if err := grpcserver.Init(dataService, conf.Server); err != nil {
		log.Fatal().Err(err)
	}
}
