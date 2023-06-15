package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig struct {
		AppVersion        string
		Port              string
		PprofPort         string
		Mode              string
		JwtSecretKey      string
		CookieName        string
		ReadTimeout       time.Duration
		WriteTimeout      time.Duration
		SSL               bool
		CtxDefaultTimeout time.Duration
		CSRF              bool
		Debug             bool
		AccessTokenMaxAge int
	} // Server config

	Port             string `mapstructure:"port"`              // Port number
	APIKey           string `mapstructure:"api_key"`           // Api key
	ConnectionString string `mapstructure:"connection_string"` // Connection string to DB
	EmailCredential  struct {
		Email    string `mapstructure:"email"`
		Password string `mapstructure:"password"`
	} `mapstructure:"email_credential"` // Email credential for Email Server
	PostgresConfig struct {
		Host     string
		Port     string
		DBName   string
		User     string
		Password string
		TimeZone string
	} `mapstructure:"postgres_config"`
	// ======= Add more configs of structure configs ===========
	// ==========================================================

}

// Populate the configuration struct
var config Config

// Init Evironment Valiables for config
func InitConfig() (*Config, error) {
	//=============== USING OS Environments ===============
	// import "os"
	// apiKey := os.Getenv("API_KEY")

	//=============== USING FLAG - default Approach by Golang==========================
	// import "flag"
	// apiKey := flag.String("<Env Name>", "<Defaut Value>", "Usage")
	// flag.Parse()

	//=============== USING Viber =========================
	// Initialize Viper
	// Set the path to look for the configurations file
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Error unmarshaling config:", err)
		return nil, err
	}
	return &config, nil
}
