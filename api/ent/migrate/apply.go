package migrate

import (
	"embed"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/tapestrylabs/fabric/api/config"
)

//go:embed migrations
var migrations embed.FS

func Apply() error {
	databaseConfig := config.NewDatabaseConfig()
	source, err := httpfs.New(http.FS(migrations), "/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("httpfs", source, databaseConfig.ConnectionString())
	if err != nil {
		return err
	}

	return m.Up()
}
