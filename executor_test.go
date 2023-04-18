package sql

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertToPositionalParameters(t *testing.T) {
	ttype := "test"
	now := time.Now()
	query, params := convertToPositionalParameters("select id, name from table where type = :type and created_at < :now", map[string]any{
		"type": ttype,
		"now":  now,
	})
	assert.Equal(t, "select id, name from table where type = $1 and created_at < $2", query)
	assert.Equal(t, 2, len(params))
	assert.Equal(t, ttype, params[0])
	assert.Equal(t, now, params[1])
}
