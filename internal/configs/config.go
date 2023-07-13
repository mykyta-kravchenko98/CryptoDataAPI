package configs

import (
	"github.com/spf13/viper"
)

type CoinMarketCapConfig struct {
	URL    string `mapstructure:"url"`
	APIKey string `mapstructure:"api-key"`
}

type ServerConfig struct {
	GRPCPort      string `mapstructure:"grpcPort"`
	RESTPort      string `mapstructure:"restPort"`
	WebSocketPort string `mapstructure:"wsPort"`
}

type KafkaConfig struct {
	CryptodataTopic string `mapstructure:"cryptodata-topic"`
	URL             string `mapstructure:"url"`
}

type Config struct {
	CoinMarketCap CoinMarketCapConfig `mapstructure:"coin-market-cap"`
	Server        ServerConfig        `mapstructure:"server"`
	Kafka         KafkaConfig         `mapstructure:"kafka"`
}

var (
	vp     *viper.Viper
	config *Config
)

func LoadConfigs(env string) (*Config, error) {
	vp = viper.New()

	vp.SetConfigType("json")
	vp.SetConfigName(env)
	vp.AddConfigPath("../configs/")
	vp.AddConfigPath("../../configs/")
	vp.AddConfigPath("configs/")

	err := vp.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return &Config{}, err
	}

	return config, err
}

func GetConfig() *Config {
	return config
}
