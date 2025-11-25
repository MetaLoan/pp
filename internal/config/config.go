package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Log        LogConfig
	JWT        JWTConfig
	Trading    TradingConfig
	Market     MarketConfig
	Settlement SettlementConfig
	Funds      FundsConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type LogConfig struct {
	Level    string
	Filename string
}

type JWTConfig struct {
	Secret string
	Expire string
}

type TradingConfig struct {
	AllowedSymbols []string
	MinOrderAmount float64
	MaxOrderAmount float64
	MaxOpenOrders  int
	DailyLossLimit float64
}

type MarketConfig struct {
	External struct {
		Provider string
		BaseURL  string
		Path     string
		Timeout  string
		APIKey   string
	}
	Breaker struct {
		FailThreshold int
		OpenSeconds   int
	}
	MockInitial      map[string]float64
	PollInterval     string
	MaxBehindSeconds int
	SymbolMap        map[string]string
}

type SettlementConfig struct {
	Interval  string
	BatchSize int
	LockKey   int64
}

type FundsConfig struct {
	Deposit struct {
		MinAmount float64
		MaxAmount float64
		DailyMax  float64
		Cooldown  string
	}
	Withdraw struct {
		MinAmount float64
		MaxAmount float64
		DailyMax  float64
		Cooldown  string
		FeeRate   float64
		MinFee    float64
	}
}

var GlobalConfig Config

func LoadConfig() {
	configName := os.Getenv("CONFIG_FILE")
	if configName == "" {
		configName = "config"
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Allow environment overrides like SERVER_PORT, DATABASE_HOST, JWT_SECRET, TRADING_MAXOPENORDERS
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
