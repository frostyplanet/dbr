package dialect

import (
	"bytes"
	"fmt"
	"time"
	"sync"
)

type mysql struct{}

var bufferPool sync.Pool

func (d mysql) QuoteIdent(s string) string {
	return quoteIdent(s, "`")
}

func (d mysql) WriteQuoteIdent(buf Buffer, s string) {
	writeQuoteIdent(buf, s, "`")
}

func (d mysql) EncodeString(s string) string {
	buf, ok := bufferPool.Get().(*bytes.Buffer)
	if ok {
		buf.Reset()
	} else {
		buf = new(bytes.Buffer)
	}
	defer bufferPool.Put(buf)

	buf.WriteRune('\'')
	// https://dev.mysql.com/doc/refman/5.7/en/string-literals.html
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 0:
			buf.WriteString(`\0`)
		case '\'':
			buf.WriteString(`\'`)
		case '"':
			buf.WriteString(`\"`)
		case '\b':
			buf.WriteString(`\b`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case 26:
			buf.WriteString(`\Z`)
		case '\\':
			buf.WriteString(`\\`)
		default:
			buf.WriteByte(s[i])
		}
	}

	buf.WriteRune('\'')
	return buf.String()
}

func (d mysql) WriteEncodeString(buf Buffer, s string) {
	buf.WriteByte('\'')
	// https://dev.mysql.com/doc/refman/5.7/en/string-literals.html
	l := len(s)
	for i := 0; i < l; i++ {
		switch s[i] {
		case 0:
			buf.WriteString(`\0`)
		case '\'':
			buf.WriteString(`\'`)
		case '"':
			buf.WriteString(`\"`)
		case '\b':
			buf.WriteString(`\b`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		case 26:
			buf.WriteString(`\Z`)
		case '\\':
			buf.WriteString(`\\`)
		default:
			buf.WriteByte(s[i])
		}
	}

	buf.WriteByte('\'')
	return
}


func (d mysql) EncodeBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func (d mysql) EncodeTime(t time.Time) string {
	return `'` + t.UTC().Format(timeFormat) + `'`
}

func (d mysql) EncodeBytes(b []byte) string {
	if len(b) == 0 {
		return fmt.Sprintf(`''`)
	}
	return fmt.Sprintf(`0x%x`, b)
}

func (d mysql) WriteEncodeBytes(buf Buffer, b []byte) {
	if len(b) == 0 {
		buf.WriteByte('\'')
		buf.WriteByte('\'')
	} else {
		fmt.Fprintf(buf, `0x%x`, b)
	}
}

func (d mysql) Placeholder(_ int) string {
	return "?"
}
