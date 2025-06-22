package config

import (
	"os"
)

type Config struct {
	AppEnv string
	*BotConfig
	*DbConfig
}

type BotConfig struct {
	Token         string
	ClientID      string
	GuildID       string
	MainChannelID string
	WelcomeEmoji  string
	GoodbyeEmoji  string
}

type DbConfig struct {
	DbHost  string
	DbUser  string
	DbPass  string
	DbName  string
	DbPort  string
	SslMode string
}

func LoadConfig() *Config {
	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),
		BotConfig: &BotConfig{
			Token:         getEnv("TOKEN", ""),
			ClientID:      getEnv("CLIENT_ID", ""),
			GuildID:       getEnv("GUILD_ID", ""),
			MainChannelID: getEnv("MAIN_CHANNEL_ID", ""),
			WelcomeEmoji:  getEnv("WELCOME_EMOJI", ""),
			GoodbyeEmoji:  getEnv("GOODBYE_EMOJI", ""),
		},
		DbConfig: &DbConfig{
			DbHost:  getEnv("DB_HOST", "localhost"),
			DbUser:  getEnv("DB_USER", "user"),
			DbPass:  getEnv("DB_PASS", ""),
			DbName:  getEnv("DB_NAME", "discord-bot"),
			DbPort:  getEnv("DB_PORT", "5432"),
			SslMode: getEnv("DB_SSL", "disable"),
		},
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
