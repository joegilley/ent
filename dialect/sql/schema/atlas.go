package schema

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
)

type atBuilder interface {
	atTable(*Table, *schema.Table)
	atTypeC(*Column, *schema.Column) error
	atUniqueC(*Table, *Column, *schema.Table, *schema.Column)
	atIncrementC(*schema.Table, *schema.Column)
	atIncrementT(*schema.Table, int64)
	atIndex(*Index, *schema.Table, *schema.Index) error
}

func (m *Migrate) aTables(ctx context.Context, b atBuilder, tables1 []*Table) ([]*schema.Table, error) {
	tables2 := make([]*schema.Table, len(tables1))
	for i, t1 := range tables1 {
		t2 := schema.NewTable(t1.Name)
		b.atTable(t1, t2)
		if m.universalID {
			r, err := m.pkRange(ctx, m.sqlDialect, t1)
			if err != nil {
				return nil, err
			}
			b.atIncrementT(t2, r)
		}
		if err := m.aColumns(b, t1, t2); err != nil {
			return nil, err
		}
		if err := m.aIndexes(b, t1, t2); err != nil {
			return nil, err
		}
		tables2[i] = t2
	}
	for i, t1 := range tables1 {
		t2 := tables2[i]
		for _, fk1 := range t1.ForeignKeys {
			fk2 := schema.NewForeignKey(fk1.Symbol).
				SetTable(t2).
				SetOnUpdate(schema.ReferenceOption(fk1.OnUpdate)).
				SetOnDelete(schema.ReferenceOption(fk1.OnDelete))
			for _, c1 := range fk1.Columns {
				c2, ok := t2.Column(c1.Name)
				if !ok {
					return nil, fmt.Errorf("unexpected fk %q column: %q", fk1.Symbol, c1.Name)
				}
				fk2.AddColumns(c2)
			}
			var refT *schema.Table
			for _, t2 := range tables2 {
				if t2.Name == fk1.RefTable.Name {
					refT = t2
					break
				}
			}
			if refT == nil {
				return nil, fmt.Errorf("unexpected fk %q ref-table: %q", fk1.Symbol, fk1.RefTable.Name)
			}
			fk2.SetRefTable(refT)
			for _, c1 := range fk1.RefColumns {
				c2, ok := refT.Column(c1.Name)
				if !ok {
					return nil, fmt.Errorf("unexpected fk %q ref-column: %q", fk1.Symbol, c1.Name)
				}
				fk2.AddRefColumns(c2)
			}
			t2.AddForeignKeys(fk2)
		}
	}
	return tables2, nil
}

func (m *Migrate) aColumns(b atBuilder, t1 *Table, t2 *schema.Table) error {
	for _, c1 := range t1.Columns {
		c2 := schema.NewColumn(c1.Name).
			SetNull(c1.Nullable)
		if c1.Collation != "" {
			c2.SetCollation(c1.Collation)
		}
		if err := b.atTypeC(c1, c2); err != nil {
			return err
		}
		if c1.Default != nil && c1.supportDefault() {
			// Has default and the database supports adding this default.
			x := fmt.Sprint(c1.Default)
			if v, ok := c1.Default.(string); ok && c1.Type != field.TypeUUID && c1.Type != field.TypeTime {
				// Escape single quote by replacing each with 2.
				x = fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
			}
			c2.SetDefault(&schema.RawExpr{X: x})
		}
		if c1.Unique {
			b.atUniqueC(t1, c1, t2, c2)
		}
		if c1.Increment {
			b.atIncrementC(t2, c2)
		}
		t2.AddColumns(c2)
	}
	return nil
}

func (m *Migrate) aIndexes(b atBuilder, t1 *Table, t2 *schema.Table) error {
	// Primary-key index.
	pk := make([]*schema.Column, 0, len(t1.PrimaryKey))
	for _, c1 := range t1.PrimaryKey {
		c2, ok := t2.Column(c1.Name)
		if !ok {
			return fmt.Errorf("unexpected primary-key column: %q", c1.Name)
		}
		pk = append(pk, c2)
	}
	t2.SetPrimaryKey(schema.NewPrimaryKey(pk...))
	// Rest of indexes.
	for _, idx1 := range t1.Indexes {
		idx2 := schema.NewIndex(idx1.Name).
			SetUnique(idx1.Unique)
		if err := b.atIndex(idx1, t2, idx2); err != nil {
			return err
		}
		t2.AddIndexes(idx2)
	}
	return nil
}

func (d *MySQL) atTable(t1 *Table, t2 *schema.Table) {
	t2.SetCharset("utf8mb4").SetCollation("utf8mb4_bin")
	if t1.Annotation == nil {
		return
	}
	if charset := t1.Annotation.Charset; charset != "" {
		t2.SetCharset(charset)
	}
	if collate := t1.Annotation.Collation; collate != "" {
		t2.SetCollation(collate)
	}
	if opts := t1.Annotation.Options; opts != "" {
		t2.AddAttrs(&mysql.CreateOptions{
			V: opts,
		})
	}
	if check := t1.Annotation.Check; check != "" {
		t2.AddChecks(&schema.Check{
			Expr: check,
		})
	}
	if checks := t1.Annotation.Checks; len(t1.Annotation.Checks) > 0 {
		names := make([]string, 0, len(checks))
		for name := range checks {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			t2.AddChecks(&schema.Check{
				Name: name,
				Expr: checks[name],
			})
		}
	}
}

func (d *MySQL) atTypeC(c1 *Column, c2 *schema.Column) error {
	if c1.SchemaType != nil && c1.SchemaType[dialect.MySQL] != "" {
		t, err := mysql.ParseType(strings.ToLower(c1.SchemaType[dialect.MySQL]))
		if err != nil {
			return err
		}
		c2.Type.Type = t
		return nil
	}
	var t schema.Type
	switch c1.Type {
	case field.TypeBool:
		t = &schema.BoolType{T: "boolean"}
	case field.TypeInt8:
		t = &schema.IntegerType{T: mysql.TypeTinyInt}
	case field.TypeUint8:
		t = &schema.IntegerType{T: mysql.TypeTinyInt, Unsigned: true}
	case field.TypeInt16:
		t = &schema.IntegerType{T: mysql.TypeSmallInt}
	case field.TypeUint16:
		t = &schema.IntegerType{T: mysql.TypeSmallInt, Unsigned: true}
	case field.TypeInt32:
		t = &schema.IntegerType{T: mysql.TypeInt}
	case field.TypeUint32:
		t = &schema.IntegerType{T: mysql.TypeInt, Unsigned: true}
	case field.TypeInt, field.TypeInt64:
		t = &schema.IntegerType{T: mysql.TypeBigInt}
	case field.TypeUint, field.TypeUint64:
		t = &schema.IntegerType{T: mysql.TypeBigInt, Unsigned: true}
	case field.TypeBytes:
		size := int64(math.MaxUint16)
		if c1.Size > 0 {
			size = c1.Size
		}
		switch {
		case size <= math.MaxUint8:
			t = &schema.BinaryType{T: mysql.TypeTinyBlob}
		case size <= math.MaxUint16:
			t = &schema.BinaryType{T: mysql.TypeBlob}
		case size < 1<<24:
			t = &schema.BinaryType{T: mysql.TypeMediumBlob}
		case size <= math.MaxUint32:
			t = &schema.BinaryType{T: mysql.TypeLongBlob}
		}
	case field.TypeJSON:
		t = &schema.JSONType{T: mysql.TypeJSON}
		if compareVersions(d.version, "5.7.8") == -1 {
			t = &schema.BinaryType{T: mysql.TypeLongBlob}
		}
	case field.TypeString:
		size := c1.Size
		if size == 0 {
			size = d.defaultSize(c1)
		}
		switch {
		case c1.typ == "tinytext", c1.typ == "text":
			t = &schema.StringType{T: c1.typ}
		case size <= math.MaxUint16:
			t = &schema.StringType{T: mysql.TypeVarchar, Size: int(size)}
		case size == 1<<24-1:
			t = &schema.StringType{T: mysql.TypeMediumText}
		default:
			t = &schema.StringType{T: mysql.TypeLongText}
		}
	case field.TypeFloat32, field.TypeFloat64:
		t = &schema.FloatType{T: c1.scanTypeOr(mysql.TypeDouble)}
	case field.TypeTime:
		t = &schema.TimeType{T: c1.scanTypeOr(mysql.TypeTimestamp)}
		// In MariaDB or in MySQL < v8.0.2, the TIMESTAMP column has both `DEFAULT CURRENT_TIMESTAMP`
		// and `ON UPDATE CURRENT_TIMESTAMP` if neither is specified explicitly. this behavior is
		// suppressed if the column is defined with a `DEFAULT` clause or with the `NULL` attribute.
		if _, maria := d.mariadb(); maria || compareVersions(d.version, "8.0.2") == -1 && c1.Default == nil {
			c2.SetNull(c1.Attr == "")
		}
	case field.TypeEnum:
		t = &schema.EnumType{T: mysql.TypeEnum, Values: c1.Enums}
	case field.TypeUUID:
		// "CHAR(X) BINARY" is treated as "CHAR(X) COLLATE latin1_bin", and in MySQL < 8,
		// and "COLLATE utf8mb4_bin" in MySQL >= 8. However we already set the table to
		t = &schema.StringType{T: mysql.TypeChar, Size: 36}
		c2.SetCollation("utf8mb4_bin")
	default:
		t, err := mysql.ParseType(c1.typ)
		if err != nil {
			return err
		}
		c2.Type.Type = t
	}
	c2.Type.Type = t
	return nil
}

func (d *MySQL) atUniqueC(t1 *Table, c1 *Column, t2 *schema.Table, c2 *schema.Column) {
	// For UNIQUE columns, MySQL create an implicit index
	// named as the column with an extra index in case the
	// name is already taken (<e.g. c>, <c_2>, <c_3>, ...).
	for _, idx := range t1.Indexes {
		// Index also defined explicitly, and will be add in atIndexes.
		if idx.Unique && d.atImplicitIndexName(idx, c1) {
			return
		}
	}
	t2.AddIndexes(schema.NewUniqueIndex(c1.Name).AddColumns(c2))
}

func (d *MySQL) atIncrementC(_ *schema.Table, c *schema.Column) {
	c.AddAttrs(&mysql.AutoIncrement{})
}

func (d *MySQL) atIncrementT(t *schema.Table, v int64) {
	t.AddAttrs(&mysql.AutoIncrement{V: v})
}

func (d *MySQL) atImplicitIndexName(idx *Index, c1 *Column) bool {
	if idx.Name == c1.Name {
		return true
	}
	if !strings.HasPrefix(idx.Name, c1.Name+"_") {
		return false
	}
	i, err := strconv.ParseInt(strings.TrimLeft(idx.Name, c1.Name+"_"), 10, 64)
	return err == nil && i > 1
}

func (d *MySQL) atIndex(idx1 *Index, t2 *schema.Table, idx2 *schema.Index) error {
	prefix := indexParts(idx1)
	for _, c1 := range idx1.Columns {
		c2, ok := t2.Column(c1.Name)
		if !ok {
			return fmt.Errorf("unexpected index %q column: %q", idx1.Name, c1.Name)
		}
		part := &schema.IndexPart{C: c2}
		if v, ok := prefix[c1.Name]; ok {
			part.AddAttrs(&mysql.SubPart{Len: int(v)})
		}
		idx2.AddParts(part)
	}
	return nil
}
