// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"github.com/jogly/ent/dialect/sql/schema"
	"github.com/jogly/ent/schema/field"
)

var (
	// CitiesColumns holds the columns for the "cities" table.
	CitiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
	}
	// CitiesTable holds the schema information for the "cities" table.
	CitiesTable = &schema.Table{
		Name:       "cities",
		Columns:    CitiesColumns,
		PrimaryKey: []*schema.Column{CitiesColumns[0]},
	}
	// StreetsColumns holds the columns for the "streets" table.
	StreetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "city_streets", Type: field.TypeInt, Nullable: true},
	}
	// StreetsTable holds the schema information for the "streets" table.
	StreetsTable = &schema.Table{
		Name:       "streets",
		Columns:    StreetsColumns,
		PrimaryKey: []*schema.Column{StreetsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "streets_cities_streets",
				Columns:    []*schema.Column{StreetsColumns[2]},
				RefColumns: []*schema.Column{CitiesColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "street_name_city_streets",
				Unique:  true,
				Columns: []*schema.Column{StreetsColumns[1], StreetsColumns[2]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		CitiesTable,
		StreetsTable,
	}
)

func init() {
	StreetsTable.ForeignKeys[0].RefTable = CitiesTable
}
