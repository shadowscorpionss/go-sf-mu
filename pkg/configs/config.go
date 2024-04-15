package config

import "os"

type Config struct {
	Censor   FConfig
	Comments FConfig
	News     FConfig
	Gateway  SConfig
}

type SConfig struct {
	AdrPort string
}

type FConfig struct {
	URLdb string
	SConfig
}

// creates new Config
func New() *Config {
	return &Config{

		Censor: FConfig{
			URLdb:   getEnv("CENSOR_DB", ""),
			SConfig: SConfig{AdrPort: getEnv("CENSOR_PORT", ":8083")},
		},
		Comments: FConfig{
			SConfig: SConfig{AdrPort: getEnv("COMMENTS_PORT", ":8082")},
			URLdb:   getEnv("COMMENTS_DB", ""),
		},
		News: FConfig{
			SConfig: SConfig{AdrPort: getEnv("NEWS_PORT", ":8081")},
			URLdb:   getEnv("NEWS_DB", ""),
		},
		Gateway: SConfig{
			AdrPort: getEnv("GATEWAY_PORT", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
