package config

import (
	"log"
	"github.com/spf13/viper"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`
	Env        string `mapstructure:"ENV"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	Port       string `mapstructure:"PORT"`
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file, continuing with environment variables")
	}
	
	viper.AutomaticEnv()
	
	// Set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "rubyone_user")
	viper.SetDefault("DB_PASSWORD", "your_secure_password_here")
	viper.SetDefault("DB_NAME", "rubyone_voice_db")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("JWT_SECRET", "your-secret-key")
	viper.SetDefault("PORT", "8080")
	
	config := &Config{}
	
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Failed to unmarshal config:", err)
	}
	
	AppConfig = config
}

func GetConfig() *Config {
	if AppConfig == nil {
		LoadConfig()
	}
	return AppConfig
} 