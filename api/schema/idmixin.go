package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/lucsky/cuid"
)

// IDMixin holds the schema definition for the IDMixin entity.
type IDMixin struct {
	mixin.Schema
}

// Fields of the IDMixin.
func (IDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Text("id").
			DefaultFunc(cuid.New).
			Immutable(),
	}
}
