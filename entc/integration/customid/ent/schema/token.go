// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package schema

import (
	"github.com/jogly/ent"
	"github.com/jogly/ent/dialect"
	"github.com/jogly/ent/entc/integration/customid/sid"
	"github.com/jogly/ent/schema/edge"
	"github.com/jogly/ent/schema/field"
)

// Token holds the schema definition for the Token entity.
type Token struct {
	ent.Schema
}

// Fields of the Token.
func (Token) Fields() []ent.Field {
	return []ent.Field{
		field.Other("id", sid.ID("")).
			SchemaType(map[string]string{
				dialect.MySQL:    "bigint",
				dialect.Postgres: "bigint",
				dialect.SQLite:   "integer",
			}).
			Default(sid.New).
			Immutable(),
		field.String("body").NotEmpty(),
	}
}

// Edges of the Token.
func (Token) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("account", Account.Type).Ref("token").Required().Unique(),
	}
}
