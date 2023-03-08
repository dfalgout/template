package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TimeMixin holds the schema definition for the TimeMixin entity.
type TimeMixin struct {
	mixin.Schema
}

// Fields of the TimeMixin.
func (TimeMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			Immutable().
			Annotations(
				entgql.OrderField("CREATED_AT"),
				entgql.Skip(entgql.SkipMutationCreateInput),
				entgql.Skip(entgql.SkipMutationUpdateInput),
			),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			SchemaType(map[string]string{
				dialect.Postgres: "timestamptz",
			}).
			Immutable().
			Annotations(
				entgql.OrderField("UPDATED_AT"),
				entgql.Skip(entgql.SkipMutationCreateInput),
				entgql.Skip(entgql.SkipMutationUpdateInput),
			),
	}
}
