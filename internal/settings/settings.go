package settings

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// App holds onto the application specific configurations.
type App struct {
	Database Database `json:"database"`
	Ping     Ping     `json:"ping"`
}

// Database holds onto the Database specific configurations.
type Database struct {
	Table  string `json:"table"`
	DSN    string `json:"dsn"`    // CBA to break down the dumb fields.
	Driver string `json:"driver"` // used for testing, mostly.
}

// Ping holds onto the Ping specific configurations.
type Ping struct {
	Locations []string `json:"locations"`
}

// Load loads the settings.
func Load() (App, error) {
	var app App

	if err := godotenv.Load(); err != nil {
		// We don't care if an .env is missing, it will be in prod.
		if !os.IsNotExist(err) {
			return app, err
		}
	}

	if err := envconfig.Process("dropout", &app); err != nil {
		return app, err
	}

	return app, nil
}
