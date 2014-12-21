package wally

import (
	"strconv"

	rdb "github.com/dancannon/gorethink"
)

func DatabaseRebuild(session *rdb.Session) {
	rdb.Db(Database).TableDrop(DocumentTable).Exec(session)
	rdb.Db(Database).TableDrop(IndexTable).Exec(session)
	rdb.Db(Database).TableCreate(DocumentTable).Exec(session)
	rdb.Db(Database).TableCreate(IndexTable).Exec(session)
	rdb.Db(Database).Table(IndexTable).IndexCreate("word").Exec(session)
}

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
