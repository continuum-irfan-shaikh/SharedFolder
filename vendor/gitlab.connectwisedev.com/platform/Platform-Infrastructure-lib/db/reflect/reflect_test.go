package reflect

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Tweet struct {
	Timeline      string
	ID            gocql.UUID  `cql:"id"`
	Text          string      `teXt`
	OriginalTweet *gocql.UUID `json:"origin"`
}

func TestStructToMap(t *testing.T) {
	//Test that if the value is not a struct we return nil, false
	m, ok := StructToMap("str")
	assert.Nil(t, m)
	assert.False(t, ok)

	tweet := Tweet{
		"t",
		gocql.TimeUUID(),
		"hello gocassa",
		nil,
	}

	m, ok = StructToMap(tweet)
	assert.True(t, ok)

	assert.Equal(t, tweet.Timeline, m["Timeline"])
	assert.Equal(t, tweet.ID, m["id"])
	assert.Equal(t, tweet.Text, m["teXt"])
	assert.Equal(t, tweet.OriginalTweet, m["OriginalTweet"])

	id := gocql.TimeUUID()
	tweet.OriginalTweet = &id
	m, _ = StructToMap(tweet)
	assert.Equal(t, tweet.OriginalTweet, m["OriginalTweet"])
}

func TestMapToStruct(t *testing.T) {
	m := make(map[string]interface{})
	assertCall := func() {
		tweet := Tweet{}
		require.NoError(t, MapToStruct(m, &tweet))
		assertEntryEqual(t, m, "Timeline", "", tweet.Timeline)
		assertEntryEqual(t, m, "id", gocql.UUID{}, tweet.ID)
		assertEntryEqual(t, m, "teXt", "", tweet.Text)
		var uuid *gocql.UUID
		assertEntryEqual(t, m, "OriginalTweet", uuid, tweet.OriginalTweet)
	}

	assertCall()
	m["Timeline"] = "timeline"
	assertCall()
	m["id"] = gocql.TimeUUID()
	assertCall()
	m["text"] = "Hello gocassa"
	assertCall()
	id := gocql.TimeUUID()
	m["OriginalTweet"] = &id
	assertCall()
}

func TestFieldsAndValues(t *testing.T) {
	var emptyUUID gocql.UUID
	id := gocql.TimeUUID()
	var nilID *gocql.UUID
	var tests = []struct {
		tweet  interface{}
		fields []string
		values []interface{}
	}{
		{
			Tweet{},
			[]string{"Timeline", "id", "teXt", "OriginalTweet"},
			[]interface{}{"", emptyUUID, "", nilID},
		},
		{
			Tweet{"timeline1", id, "hello gocassa", &id},
			[]string{"Timeline", "id", "teXt", "OriginalTweet"},
			[]interface{}{"timeline1", id, "hello gocassa", &id},
		},
		{
			nil,
			[]string{},
			[]interface{}{},
		},
	}
	for _, test := range tests {
		fields, values, _ := FieldsAndValues(test.tweet)
		assertFieldsEqual(t, test.fields, fields)
		assertValuesEqual(t, test.values, values)
	}
}

func assertEntryEqual(t *testing.T, m map[string]interface{}, key string, empty, actual interface{}) {
	expected, ok := m[key]
	if ok {
		if expected != actual {
			t.Errorf("Expected %q to be %s but got %s", key, expected, actual)
		}
	} else {
		if actual != empty {
			t.Errorf("Expected %q to be empty but got %s", key, actual)
		}
	}
}

func assertFieldsEqual(t *testing.T, a, b []string) {
	if len(a) != len(b) {
		t.Errorf("expected fields %v but got %v", a, b)
		return
	}

	for i := range a {
		if a[i] != b[i] {
			t.Errorf("expected fields %v but got %v", a, b)
		}
	}
}

func assertValuesEqual(t *testing.T, a, b []interface{}) {
	if len(a) != len(b) {
		t.Errorf("expected values %v but got %v different length", a, b)
		return
	}

	for i := range a {
		if a[i] != b[i] {
			t.Errorf("expected values %v but got %v a[i] = %v and b[i] = %v", a, b, a[i], b[i])
			return
		}
	}
}
