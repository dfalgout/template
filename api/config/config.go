package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()

	// Set deafults here instead of loading .env since we aren't using .env in cloud deployment envs
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5436")
	viper.SetDefault("DB_NAME", "fabric")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASSWORD", "password")

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DEBUG", false)
	viper.SetDefault("AUTO_MIGRATIONS", false)
	viper.SetDefault("ENVIRONMENT", "local")
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     viper.GetString("DB_HOST"),
		Port:     viper.GetString("DB_PORT"),
		Name:     viper.GetString("DB_NAME"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
	}
}

func (dc *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@/%s?host=%s&port=%s&sslmode=disable", dc.User, dc.Password, dc.Name, dc.Host, dc.Port)
}

type ServerConfig struct {
	Port           string
	Debug          bool
	AutoMigrations bool
	Environment    string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:           viper.GetString("PORT"),
		Debug:          viper.GetBool("DEBUG"),
		AutoMigrations: viper.GetBool("AUTO_MIGRATIONS"),
		Environment:    viper.GetString("ENVIRONMENT"),
	}
}

func (sc *ServerConfig) PortString() string {
	return fmt.Sprintf(":%s", sc.Port)
}

// Left in as a resource for future wiring if redirects are required in backend

// func (sc ServerConfig) AssetsEndpoint() string {
// 	subdomain := fmt.Sprintf("%s-assets", sc.Environment)
// 	if sc.Environment == "prod" {
// 		subdomain = "assets"
// 	}
// 	return fmt.Sprintf("https://%s.yummi.ninja", subdomain)
// }

// func (sc ServerConfig) AssetsBucket() string {
// 	subdomain := fmt.Sprintf("yn-assets-%s", sc.Environment)
// 	if sc.Environment == "prod" {
// 		subdomain = "yn-assets"
// 	}

// 	return subdomain
// }

func (sc ServerConfig) BaseEndpoint() string {
	protocol := "http"
	domain := fmt.Sprintf("%s:%s", "localhost", sc.Port)
	// if sc.Environment == "dev" {
	// 	protocol = "https"
	// 	domain = "dev-api.yummi.ninja"
	// } else if sc.Environment == "prod" {
	// 	protocol = "https"
	// 	domain = "api.yummi.ninja"
	// }

	return fmt.Sprintf("%s://%s", protocol, domain)
}
