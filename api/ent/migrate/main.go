//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	"github.com/tapestrylabs/fabric/api/config"
	"github.com/tapestrylabs/fabric/api/ent/migrate"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"
)

func main() {
	databaseConfig := config.NewDatabaseConfig()
	ctx := context.Background()
	// Create a local migration directory able to understand Atlas migration file format for replay.
	dir, err := sqltool.NewGolangMigrateDir("ent/migrate/migrations")
	if err != nil {
		log.Fatalf("failed creating atlas migration directory: %v", err)
	}
	// Migrate diff options.
	opts := []schema.MigrateOption{
		schema.WithDir(dir),                         // provide migration directory
		schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
		schema.WithDialect(dialect.Postgres),        // Ent dialect to use
	}
	if len(os.Args) != 2 {
		log.Fatalln("migration name is required. Use: 'go run ent/migrate/main.go <name>'")
	}
	err = migrate.NamedDiff(ctx, databaseConfig.ConnectionString(), os.Args[1], opts...)
	if err != nil {
		log.Fatalf("failed generating migration file: %v", err)
	}
}
