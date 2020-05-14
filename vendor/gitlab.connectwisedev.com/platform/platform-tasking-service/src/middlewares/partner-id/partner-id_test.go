package partnerID

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	name      string
	passed    bool
	partnerID string
	url       string
	mustPass  bool
}

func (t *testCase) next(_ http.ResponseWriter, r *http.Request) {
	gotPartnerID := r.Context().Value(PartnerIDKey)

	if gotPartnerID == t.partnerID {
		t.passed = true
	}
}

func TestMiddleware(t *testing.T) {
	recorder := httptest.NewRecorder()

	testCases := []testCase{
		{
			name:      "good_partner_ID_string",
			partnerID: "partnerID",
			url:       "http://localhost/partners",
			mustPass:  true,
		},
		{
			name:      "good_partner_ID_number",
			partnerID: "12345",
			url:       "http://localhost/partners",
			mustPass:  true,
		},
		{
			name:      "bad_url",
			partnerID: "partnerID",
			url:       "http://localhost/parts",
			mustPass:  false,
		},
		{
			name:      "bad_partner_ID",
			partnerID: "",
			url:       "http://localhost/parts",
			mustPass:  false,
		},
	}
	fmt.Println(pattern)
	for _, testCase := range testCases {
		test := testCase
		t.Run(test.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/%s/", test.url, test.partnerID)
			request := httptest.NewRequest("GET", url, nil)

			Middleware(recorder, request.WithContext(context.Background()), test.next)
			if test.mustPass != test.passed {
				t.Fatalf("want mustPass=%t but it did not", test.mustPass)
			}
		})
	}
}
