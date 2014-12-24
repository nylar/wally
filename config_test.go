package wally

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConf = `
database:
  host: localhost:28015
  name: wally

tables:
  document_table: documents
  index_table: indexes
`

func TestConfig_LoadConfig(t *testing.T) {
	data := []byte(validConf)
	conf, err := LoadConfig(data)

	assert.NoError(t, err)
	assert.IsType(t, conf, new(Config))

	assert.Equal(t, conf.Database.Host, "localhost:28015")
	assert.Equal(t, conf.Database.Name, "wally")
	assert.Equal(t, conf.Tables.DocumentTable, "documents")
	assert.Equal(t, conf.Tables.IndexTable, "indexes")
}

func TestConfig_LoadConfigBadYAML(t *testing.T) {
	badYAML := `database!`

	data := []byte(badYAML)
	conf, err := LoadConfig(data)
	assert.Error(t, err)
	assert.Nil(t, conf)
}
