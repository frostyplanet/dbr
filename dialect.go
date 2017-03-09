package dbr

import (
	"time"
	"github.com/gocraft/dbr/dialect"
)

// Dialect abstracts database differences
type Dialect interface {
	QuoteIdent(id string) string
	WriteQuoteIdent(buf dialect.Buffer, id string)

	EncodeString(s string) string
	WriteEncodeString(buf dialect.Buffer, s string)
	EncodeBool(b bool) string
	EncodeTime(t time.Time) string
	EncodeBytes(b []byte) string
	WriteEncodeBytes(buf dialect.Buffer, b []byte)

	Placeholder(n int) string
}
