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
	// GroupsColumns holds the columns for the "groups" table.
	GroupsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "group_admin", Type: field.TypeInt, Nullable: true},
	}
	// GroupsTable holds the schema information for the "groups" table.
	GroupsTable = &schema.Table{
		Name:       "groups",
		Columns:    GroupsColumns,
		PrimaryKey: []*schema.Column{GroupsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "groups_users_admin",
				Columns:    []*schema.Column{GroupsColumns[2]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// PetsColumns holds the columns for the "pets" table.
	PetsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "name", Type: field.TypeString},
		{Name: "user_pets", Type: field.TypeInt, Nullable: true},
	}
	// PetsTable holds the schema information for the "pets" table.
	PetsTable = &schema.Table{
		Name:       "pets",
		Columns:    PetsColumns,
		PrimaryKey: []*schema.Column{PetsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "pets_users_pets",
				Columns:    []*schema.Column{PetsColumns[2]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "age", Type: field.TypeInt},
		{Name: "name", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// GroupUsersColumns holds the columns for the "group_users" table.
	GroupUsersColumns = []*schema.Column{
		{Name: "group_id", Type: field.TypeInt},
		{Name: "user_id", Type: field.TypeInt},
	}
	// GroupUsersTable holds the schema information for the "group_users" table.
	GroupUsersTable = &schema.Table{
		Name:       "group_users",
		Columns:    GroupUsersColumns,
		PrimaryKey: []*schema.Column{GroupUsersColumns[0], GroupUsersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "group_users_group_id",
				Columns:    []*schema.Column{GroupUsersColumns[0]},
				RefColumns: []*schema.Column{GroupsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "group_users_user_id",
				Columns:    []*schema.Column{GroupUsersColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// PetFriendsColumns holds the columns for the "pet_friends" table.
	PetFriendsColumns = []*schema.Column{
		{Name: "pet_id", Type: field.TypeInt},
		{Name: "friend_id", Type: field.TypeInt},
	}
	// PetFriendsTable holds the schema information for the "pet_friends" table.
	PetFriendsTable = &schema.Table{
		Name:       "pet_friends",
		Columns:    PetFriendsColumns,
		PrimaryKey: []*schema.Column{PetFriendsColumns[0], PetFriendsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "pet_friends_pet_id",
				Columns:    []*schema.Column{PetFriendsColumns[0]},
				RefColumns: []*schema.Column{PetsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "pet_friends_friend_id",
				Columns:    []*schema.Column{PetFriendsColumns[1]},
				RefColumns: []*schema.Column{PetsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UserFriendsColumns holds the columns for the "user_friends" table.
	UserFriendsColumns = []*schema.Column{
		{Name: "user_id", Type: field.TypeInt},
		{Name: "friend_id", Type: field.TypeInt},
	}
	// UserFriendsTable holds the schema information for the "user_friends" table.
	UserFriendsTable = &schema.Table{
		Name:       "user_friends",
		Columns:    UserFriendsColumns,
		PrimaryKey: []*schema.Column{UserFriendsColumns[0], UserFriendsColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "user_friends_user_id",
				Columns:    []*schema.Column{UserFriendsColumns[0]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "user_friends_friend_id",
				Columns:    []*schema.Column{UserFriendsColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		GroupsTable,
		PetsTable,
		UsersTable,
		GroupUsersTable,
		PetFriendsTable,
		UserFriendsTable,
	}
)

func init() {
	GroupsTable.ForeignKeys[0].RefTable = UsersTable
	PetsTable.ForeignKeys[0].RefTable = UsersTable
	GroupUsersTable.ForeignKeys[0].RefTable = GroupsTable
	GroupUsersTable.ForeignKeys[1].RefTable = UsersTable
	PetFriendsTable.ForeignKeys[0].RefTable = PetsTable
	PetFriendsTable.ForeignKeys[1].RefTable = PetsTable
	UserFriendsTable.ForeignKeys[0].RefTable = UsersTable
	UserFriendsTable.ForeignKeys[1].RefTable = UsersTable
}
