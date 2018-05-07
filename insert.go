package dbr

import (
	"bytes"
	"reflect"
)

// InsertStmt builds `INSERT INTO ...`
type InsertStmt struct {
	raw

	Table  string
	Column []string
	Value  [][]interface{}
	Ignore bool
}

// Build builds `INSERT INTO ...` in dialect
func (b *InsertStmt) Build(d Dialect, buf Buffer) error {
	if b.raw.Query != "" {
		return b.raw.Build(d, buf)
	}

	if b.Table == "" {
		return ErrTableNotSpecified
	}

	if len(b.Column) == 0 {
		return ErrColumnNotSpecified
	}
	if b.Ignore {
		buf.WriteString("INSERT IGNORE INTO ")
	} else {
		buf.WriteString("INSERT INTO ")
	}
	d.WriteQuoteIdent(buf, b.Table)

	placeholderBuf, ok := stdBufferPool.Get().(*bytes.Buffer)
	if ok {
		placeholderBuf.Reset()
	} else {
		placeholderBuf = new(bytes.Buffer)
	}
	defer stdBufferPool.Put(placeholderBuf)
	placeholderBuf.WriteString("(")
	buf.WriteString(" (")
	for i, col := range b.Column {
		if i > 0 {
			buf.WriteString(",")
			placeholderBuf.WriteString(",")
		}
		d.WriteQuoteIdent(buf, col)
		placeholderBuf.WriteString(placeholder)
	}
	buf.WriteString(") VALUES ")
	placeholderBuf.WriteString(")")
	placeholderStr := placeholderBuf.String()

	for i, tuple := range b.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(placeholderStr)

		buf.WriteValue(tuple...)
	}

	return nil
}

// InsertInto creates an InsertStmt
func InsertInto(table string) *InsertStmt {
	return &InsertStmt{
		Table: table,
	}
}

func InsertIgnoreInto(table string) *InsertStmt {
	return &InsertStmt{
		Table: table,
		Ignore: true,
	}
}

// InsertBySql creates an InsertStmt from raw query
func InsertBySql(query string, value ...interface{}) *InsertStmt {
	return &InsertStmt{
		raw: raw{
			Query: query,
			Value: value,
		},
	}
}

// Columns adds columns
func (b *InsertStmt) Columns(column ...string) *InsertStmt {
	b.Column = column
	return b
}

// Values adds a tuple for columns
func (b *InsertStmt) Values(value ...interface{}) *InsertStmt {
	b.Value = append(b.Value, value)
	return b
}

// Record adds a tuple for columns from a struct
func (b *InsertStmt) Record(structValue interface{}) *InsertStmt {
	v := reflect.Indirect(reflect.ValueOf(structValue))

	if v.Kind() == reflect.Struct {
		var value []interface{}
		m := structMap(v)
		for _, key := range b.Column {
			if val, ok := m[key]; ok {
				value = append(value, val.Interface())
			} else {
				value = append(value, nil)
			}
		}
		b.Values(value...)
	}
	return b
}
