package dialect

import (
	"io"
	"strings"
)

var (
	// MySQL dialect
	MySQL = mysql{}
	// PostgreSQL dialect
	PostgreSQL = postgreSQL{}
	// SQLite3 dialect
	SQLite3 = sqlite3{}
)

const (
	timeFormat = "2006-01-02 15:04:05.000000"
)

type Buffer interface {
	io.Writer
	io.ByteWriter
	WriteString(s string) (n int, err error)
}


func quoteIdent(s, quote string) string {
	part := strings.SplitN(s, ".", 2)
	if len(part) == 2 {
		return quoteIdent(part[0], quote) + "." + quoteIdent(part[1], quote)
	}
	return quote + s + quote
}

func writeQuoteIdent(buf Buffer, s, quote string) {
	p := strings.IndexByte(s, '.')
	if p == -1 {
		buf.WriteString(quote)
		buf.WriteString(s)
		buf.WriteString(quote)
	} else {
		writeQuoteIdent(buf, s[0:p], quote)
		buf.WriteByte('.')
		writeQuoteIdent(buf, s[p+1:len(s)], quote)
	}
}
