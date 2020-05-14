package transactionID

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

var uuidPattern = "^([a-z0-9]){8}-([a-z0-9]){4}-([a-z0-9]){4}-([a-z0-9]){4}-([a-z0-9]){12}$"

type testCase struct {
	name     string
	passed   bool
	mustPass bool
}

func (t *testCase) next(_ http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(uuidPattern)
	gotTransactionIDKey := r.Context().Value(Key)

	result := re.MatchString(gotTransactionIDKey.(string))
	if result {
		t.passed = true
	}
}

func TestFromContext(t *testing.T) {
	want := "value"
	ctx := context.WithValue(context.Background(), Key, want)

	got := FromContext(ctx)
	if got != want {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}

func TestNewContextWithID(t *testing.T) {
	got := NewContext()
	if got == nil {
		t.Fatalf("want: %v, got: %v", "not nil", got)
	}

	if id := FromContext(got); id == "" {
		t.Fatalf("want: %v, got: %v", "not empty", got)
	}
}

func TestGenerateTransactionID(t *testing.T) {
	re := regexp.MustCompile(uuidPattern)
	got := generateTransactionID()

	result := re.MatchString(got)
	if !result {
		t.Fatalf("string did not match uuid pattern, got: %q", got)
	}
}

func TestNewContextWithTransactionID(t *testing.T) {
	re := regexp.MustCompile(uuidPattern)
	ctx := context.Background()
	request := httptest.NewRequest("GET", "http://url", nil)

	got := newContextTransactionIDWithReq(ctx, request).Value(Key)
	result := re.MatchString(got.(string))
	if !result {
		t.Fatalf("newContextTransactionIDWithReq(): string did not match uuid pattern, got: %q", got)
	}
}

func TestMiddleware(t *testing.T) {
	recorder := httptest.NewRecorder()

	tt := []testCase{
		{
			name:     "good_case",
			mustPass: true,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://localhost", nil)

			Middleware(recorder, request.WithContext(context.Background()), test.next)
			if test.mustPass != test.passed {
				t.Fatalf("want mustPass=%t but it did not", test.mustPass)
			}
		})
	}
}
