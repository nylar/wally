package wally

import (
	"strconv"
	"strings"

	rdb "github.com/dancannon/gorethink"
	"github.com/fatih/color"
)

var (
	// Success coloured text
	Success = color.New(color.FgGreen)
	// Info coloured text
	Info = color.New(color.FgBlue)
	// Warning coloured text
	Warning = color.New(color.FgYellow)
	// Std coloured text
	Std = color.New(color.FgMagenta)
)

// DatabaseRebuild resets the database to an empty state, it also sets the
// secondary index for the index table.
func DatabaseRebuild(session *rdb.Session) {
	rdb.Db(Conf.Database.Name).TableDrop(Conf.Tables.DocumentTable).Exec(session)
	rdb.Db(Conf.Database.Name).TableDrop(Conf.Tables.IndexTable).Exec(session)
	rdb.Db(Conf.Database.Name).TableCreate(Conf.Tables.DocumentTable).Exec(session)
	rdb.Db(Conf.Database.Name).TableCreate(Conf.Tables.IndexTable).Exec(session)
	rdb.Db(Conf.Database.Name).Table(Conf.Tables.IndexTable).IndexCreate("word").Exec(session)
}

// ToString converts an interface{} to a string, a string, byte slice or integer
// is an accepted value and converted as such, anything else returns an empty string.
func ToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case []byte:
		return string(v.([]byte))
	case int:
		return strconv.Itoa(v.(int))
	default:
		return ""
	}
}

// TruncateText returns a truncated version of a string with a suffix appended.
func TruncateText(text, suffix string, upperBound int) string {
	if len(text) > upperBound {
		truncated := text[:upperBound]
		nextSpace := strings.LastIndex(truncated, " ")
		if nextSpace > 0 {
			truncated = text[:nextSpace]
		}
		return truncated + suffix
	}
	return text
}
