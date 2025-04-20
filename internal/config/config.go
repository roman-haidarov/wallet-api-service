package config

import (
	"os"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel    int    `mapstructure:"log_level"`
	AppBaseURL  string `mapstructure:"app_base_url"`
	AppBasePath string `mapstructure:"app_base_path"`
	Server      Server `mapstructure:"server"`
	Kafka       Kafka  `mapstructure:"kafka"`
}

type Server struct {
	ListenPort int `mapstructure:"listen_port"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
	UI      string   `mapstructure:"ui"`
	Topic   string   `mapstructure:"topic"`
}

func Load(path string) (*Config, error) {
	setDefaults()

	var c Config

	if path != "" {
		_, err := os.Lstat(path)
		if err != nil {
			return nil, err
		}
		viper.SetConfigFile(path)

		err = viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	envVars(&c)

	return &c, nil
}

func envVars(c *Config) {
	viper.AutomaticEnv()

	if val := viper.GetString("LOG_LEVEL"); val != "" {
		if intVal := viper.GetInt("LOG_LEVEL"); intVal >= 0 {
			c.LogLevel = intVal
		}
	}

	// // Jaeger
	// if val := viper.GetString("TELEMETRY_ENDPOINT"); val != "" {
	// 	c.TelemetryEndpoint = val
	// }

	// // Swagger
	// if val := viper.GetString("APP_BASE_URL"); val != "" {
	// 	c.AppBaseURL = val
	// }
	// if val := viper.GetString("APP_BASE_PATH"); val != "" {
	// 	c.AppBasePath = val
	// }
}

func setDefaults() {
	viper.SetDefault("log_level", 1)
	viper.SetDefault("server.listen_port", 8080)
	viper.SetDefault("kafka.brokers", []string{"localhost:9091", "localhost:9092", "localhost:9093"})
	viper.SetDefault("kafka.ui", "localhost:9020")
	viper.SetDefault("kafka.topic", "wallet-topups")
}
