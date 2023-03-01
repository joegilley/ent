// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package schema

import (
	"github.com/jogly/ent"
	"github.com/jogly/ent/entc/integration/customid/sid"
	"github.com/jogly/ent/schema/edge"
	"github.com/jogly/ent/schema/field"
)

// IntSID holds the schema definition for the IntSID entity.
type IntSID struct {
	ent.Schema
}

// Fields of the IntSid.
func (IntSID) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			GoType(sid.New()).
			Unique().
			Immutable(),
	}
}

// Edges of the IntSid.
func (IntSID) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("parent", IntSID.Type).
			Unique(),
		edge.From("children", IntSID.Type).
			Ref("parent"),
	}
}
