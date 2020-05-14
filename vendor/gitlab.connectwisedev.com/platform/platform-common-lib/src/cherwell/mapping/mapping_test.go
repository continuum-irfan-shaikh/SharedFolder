package mapping

import (
	"testing"
	"time"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/cherwell"
	"github.com/stretchr/testify/assert"
)

type testMapObject struct {
	Foo      string    `cherwell:"Foo"`
	Bar      int       `cherwell:"Bar"`
	Baz      bool      `cherwell:"Baz"`
	Created  time.Time `cherwell:"CreatedAt"`
	Comments string    `cherwell:"Notes,isHtml"`
	PtrVal   *string   `cherwell:"PointerVal"`
}

type testMapObject2 struct {
	Foo      string    `cherwell:"Foo"`
	Bar      int       `cherwell:"Bar"`
	Baz      bool      `cherwell:"Baz"`
	Created  time.Time `cherwell:"CreatedAt"`
	Comments string    `cherwell:"Notes,isHtml"`
	PtrVal   *string   `cherwell:"PointerVal,omitempty"`
}

type badMapObject struct {
	Bar []string `cherwell:"bar"`
}

var testPtrVal = "test"

func TestMapToBusinessObject(t *testing.T) {
	now := time.Now()
	cases := map[string]struct {
		source   interface{}
		dest     *cherwell.BusinessObject
		expected cherwell.BusinessObject
		error    string
	}{
		"should map valid BO": {
			source: &testMapObject{
				Foo:      "foo",
				Bar:      10,
				Baz:      true,
				Created:  now,
				Comments: `<span>test</span>`,
				PtrVal:   &testPtrVal,
			},
			dest: &cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Name: "Foo"},
					{Name: "Bar"},
					{Name: "Baz"},
					{Name: "CreatedAt"},
					{Name: "Notes"},
					{Name: "PointerVal"},
				},
			},
			expected: cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Dirty: true, Name: "Foo", Value: "foo"},
					{Dirty: true, Name: "Bar", Value: "10"},
					{Dirty: true, Name: "Baz", Value: "true"},
					{Dirty: true, Name: "CreatedAt", Value: now.Format(DefaultDateFormat)},
					{Dirty: true, Name: "Notes", HTML: `<span>test</span>`, Value: `<span>test</span>`},
					{Dirty: true, Name: "PointerVal", Value: testPtrVal},
				},
			},
		},
		"should ignore nil field": {
			source: &testMapObject{
				Foo:      "foo",
				Bar:      10,
				Baz:      true,
				Created:  now,
				Comments: `<span>test</span>`,
			},
			dest: &cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Name: "Foo"},
					{Name: "Bar"},
					{Name: "Baz"},
					{Name: "CreatedAt"},
					{Name: "Notes"},
					{Name: "PointerVal"},
				},
			},
			expected: cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Dirty: true, Name: "Foo", Value: "foo"},
					{Dirty: true, Name: "Bar", Value: "10"},
					{Dirty: true, Name: "Baz", Value: "true"},
					{Dirty: true, Name: "CreatedAt", Value: now.Format(DefaultDateFormat)},
					{Dirty: true, Name: "Notes", HTML: `<span>test</span>`, Value: `<span>test</span>`},
					{Dirty: false, Name: "PointerVal", Value: ""},
				},
			},
		},
		"should stop if field in tag doesn't exists": {
			error: "failed to map field 'Bar', field 'Bar' doesn't exists in business object #32",
			source: testMapObject{
				Foo: "foo",
				Bar: 10,
				Baz: true,
			},
			dest: &cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Name: "Foo"},
				},
			},
		},
		"should ignore an orphan field if 'omitempty' flag provided": {
			source: &testMapObject2{
				Foo:      "foo",
				Bar:      10,
				Baz:      true,
				Created:  now,
				Comments: `<span>test</span>`,
				PtrVal:   &testPtrVal,
			},
			dest: &cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Name: "Foo"},
					{Name: "Bar"},
					{Name: "Baz"},
					{Name: "CreatedAt"},
					{Name: "Notes"},
				},
			},
			expected: cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Dirty: true, Name: "Foo", Value: "foo"},
					{Dirty: true, Name: "Bar", Value: "10"},
					{Dirty: true, Name: "Baz", Value: "true"},
					{Dirty: true, Name: "CreatedAt", Value: now.Format(DefaultDateFormat)},
					{Dirty: true, Name: "Notes", HTML: `<span>test</span>`, Value: `<span>test</span>`},
				},
			},
		},
		"should return error if struct field is unsupported": {
			error:  "failed to map field 'Bar', cannot convert value to string, unsupported type []string",
			source: badMapObject{Bar: make([]string, 0)},
			dest: &cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Name: "bar"},
				},
			},
		},
		"should stop if object is nil": {
			error:  "nil value passed",
			source: nil,
		},
		"should stop if BO is nil": {
			error:  "output business object value is nil",
			source: &testMapObject{},
			dest:   nil,
		},
		"should check if mapping field destination is defined": {
			error: "failed to map field 'Foo', mapping destination field is missing",
			source: struct {
				Foo string `cherwell:""`
			}{Foo: "bar"},
			dest: &cherwell.BusinessObject{
				BusinessObjectInfo: cherwell.BusinessObjectInfo{
					ID: "32",
				},
				Fields: []cherwell.FieldTemplateItem{
					{Name: "bar"},
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatal(r)
				}
			}()

			err := MapToBusinessObject(c.source, c.dest)
			if c.error != "" {
				assert.EqualError(t, err, c.error)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, c.expected, *c.dest)
		})
	}
}
