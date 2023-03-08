package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/debug"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/tapestrylabs/fabric/api/config"
	"github.com/tapestrylabs/fabric/api/ent"
	"github.com/tapestrylabs/fabric/api/ent/migrate"
	"github.com/tapestrylabs/fabric/api/graph"
	"github.com/vektah/gqlparser/gqlerror"

	"go.uber.org/zap"

	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func main() {
	databaseConfig := config.NewDatabaseConfig()
	serverConfig := config.NewServerConfig()

	var (
		logger *zap.Logger
		err    error
	)
	if serverConfig.Environment == "prod" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		logger.Fatal("can't initialize zap logger: %v", zap.Error(err))
	}
	defer logger.Sync()

	dbClient := Open(databaseConfig.ConnectionString())
	if serverConfig.Debug {
		dbClient = dbClient.Debug()
	}
	defer dbClient.Close()

	if serverConfig.AutoMigrations {
		// auto migrations
		ctx := context.Background()
		err = dbClient.Schema.Create(ctx)
		if err != nil {
			logger.Fatal("auto migrations", zap.Error(err))
		}
		logger.Info("auto migrations complete")
	} else {
		err = migrate.Apply()
		if err != nil && err.Error() != "no change" {
			logger.Fatal("running migrations", zap.Error(err))
		}
		logger.Info("file migrations complete")
	}

	router := chi.NewRouter()

	AllowedOrigins := []string{} //[]string{"https://*.fabric.io"}
	if serverConfig.Environment == "local" {
		AllowedOrigins = []string{"*"}
	}

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   AllowedOrigins,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            serverConfig.Debug,
	}).Handler)

	server := handler.NewDefaultServer(graph.NewSchema(dbClient, logger, serverConfig))
	server.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logger.Error("fatal error encountered",
			zap.Any("error", err),
		)
		return gqlerror.Errorf("internal server error")
	})

	// Setup graphql mutation transactions
	server.Use(entgql.Transactioner{TxOpener: dbClient})
	if serverConfig.Debug {
		server.Use(&debug.Tracer{})
	}

	// Health check
	router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	// Graphql routes
	if serverConfig.Environment == "local" {
		router.Handle("/", playground.Handler("Fabric", "/graphql"))
	}
	router.Handle("/graphql", server)

	logger.Info("listening on", zap.String("address", serverConfig.PortString()))
	if err := http.ListenAndServe(serverConfig.PortString(), router); err != nil {
		logger.Error("http server terminated", zap.Error(err))
	}
}
