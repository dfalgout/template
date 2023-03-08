//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithSchemaGenerator(),
		entgql.WithWhereInputs(true),
		entgql.WithSchemaPath("../gql/ent.graphql"),
		entgql.WithConfigPath("../gqlgen.yml"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	features := []gen.Feature{
		gen.FeatureVersionedMigration,
	}

	if err := entc.Generate("../schema", &gen.Config{
		Templates: entgql.AllTemplates,
		Schema:    "api/schema",
		Target:    "../ent",
		Package:   "github.com/tapestrylabs/fabric/api/ent",
		Features:  features,
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
