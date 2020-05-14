package cherwell

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBOTemplate_NewBusinessObject(t *testing.T) {
	expected := BusinessObject{
		BusinessObjectInfo: BusinessObjectInfo{
			ID: "1",
		},
		Fields: []FieldTemplateItem{
			{
				FieldID: "11",
				Name:    "Foo",
			},
		},
	}

	tpl := BOTemplate{
		ID:     expected.ID,
		Fields: expected.Fields,
	}

	got := tpl.NewBusinessObject()
	assert.Equal(t, expected, got)
}

func TestClient_GetBusinessObjectTemplate(t *testing.T) {
	cases := map[string]struct {
		query        BOTemplateQuery
		responseBody interface{}
		responseCode int
		want         BOTemplate
		wantError    string
	}{
		"should return BO template": {
			responseCode: http.StatusOK,
			responseBody: BOTemplate{
				Fields: []FieldTemplateItem{
					{
						Name: "foo",
					},
				},
			},
			want: BOTemplate{
				ID: "testID",
				Fields: []FieldTemplateItem{
					{
						Name: "foo",
					},
				},
			},
			query: BOTemplateQuery{
				ID:         "testID",
				FieldNames: []string{"foo", "bar"},
				FieldIDs:   []string{"1", "2", "3"},
				IncludeAll: true,
			},
		},
		"should populate errors from cherwell": {
			responseCode: http.StatusInternalServerError,
			responseBody: ErrorData{
				HasError:     true,
				ErrorMessage: "Something broken",
			},
			wantError: "Something broken",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			defer catchPanics(t)

			shouldFail := c.responseCode >= 300
			server, mux := newTestServer()
			defer server.Close()

			resp, err := json.Marshal(c.responseBody)
			assert.NoError(t, err, "Can not prepare mock response: %v", err)

			mux.Handle(boTemplateURL, newMockHandler(http.MethodPost, boTemplateURL, string(resp), c.responseCode))
			client, err := NewClient(Config{Host: server.URL}, getWebClient())
			assert.NoError(t, err)

			got, err := client.GetBusinessObjectTemplate(c.query)
			if err != nil {
				if !shouldFail {
					t.Fatalf("unexpected error from Client: %s", err)
				}

				assert.EqualError(t, err, c.wantError)
				return
			}

			assert.Equal(t, c.want, *got)
		})
	}
}
