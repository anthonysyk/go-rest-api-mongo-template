package jsonflex

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt_UnmarshalJSON(t *testing.T) {
	input := `[{"value": "42"}, {"value": 11}]`

	type Foo struct {
		Value Int `json:"value"`
	}

	var values []Foo
	err := json.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)
	assert.Len(t, values, 2)
	assert.Equal(t, int(values[0].Value), 42)
	assert.Equal(t, int(values[1].Value), 11)
}
